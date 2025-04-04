from fastapi import FastAPI, Response, status, HTTPException
from googletrans import Translator
import fasttext
import time
import numpy as np
from typing import List
from pydantic import BaseModel
from iso639 import iter_langs

from .utils import clear_text


app = FastAPI()
translator = Translator()
fasttext_detector = fasttext.load_model("lid.176.bin")

@app.get("/translate")
async def translate(text: str, dest: str, src: str, response: Response = None):
    # Return 204 if source language matches destination
    if src == dest:
        response.status_code = status.HTTP_204_NO_CONTENT
        return

    # Measure translate time
    translate_start = time.time()
    try:
        translations = await translator.translate([text.strip()], dest=dest, src=src)
        if not translations:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Google services unavailable: translation failed"
            )
    except Exception:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Google services unavailable: translation failed"
        )
    translate_time = time.time() - translate_start

    print(f"Translation completed: Text: {text} | Result: {[t.text for t in translations]} | Translation time: {translate_time:.3f}s")
    
    return {
        "source_language": src,
        "source_text": text,
        "translated_text": [t.text for t in translations],
        "destination_language": dest
    }

@app.get("/detect")
async def detect(text: str):
    try:
        cleaned_text = clear_text(text.strip()).strip()
        if not cleaned_text:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Text contains no detectable content after cleaning"
            )
            
        # FastText returns tuple of (labels, probabilities)
        labels, probabilities = fasttext_detector.predict(cleaned_text, k=3)
        probabilities = np.asarray(probabilities).flatten()
        
        if not labels or len(labels) == 0:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Language detection failed: no languages detected"
            )
            
        detected_langs = [label.replace('__label__', '') for label in labels]
        
        results = [
            {"language": lang, "probability": float(prob)} 
            for lang, prob in zip(detected_langs, probabilities)
        ]
        
        print(f"Detection completed: Original text: {text} | Cleaned text: {cleaned_text} | Results: {results}")
        
        return {
            "text": text,
            "cleaned_text": cleaned_text,
            "detected_languages": results
        }
        
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Language detection failed: {str(e)}"
        )

class ResponseLang(BaseModel):
    iso_639_1: str
    name: str

LANGUAGE_MAP = {
    lang.pt1: {
        "name": lang.name,
    }
    for lang in iter_langs()
    if lang.pt1  # Only include languages with ISO 639-1 codes
}

@app.get("/detect/languages", response_model=List[ResponseLang])
async def get_supported_languages():
    # Get all labels from the model
    labels = (fasttext_detector.get_labels())
    
    # Process labels and create response
    supported_languages = []
    for label in labels:
        # Remove '__label__' prefix from fasttext labels
        lang_code = label.replace('__label__', '')

        # Get language info from our map, or use code as fallback
        lang_info = LANGUAGE_MAP.get(lang_code, {
            "name": lang_code,
            "native_name": lang_code
        })

        supported_languages.append(ResponseLang(
            iso_639_1=lang_code,
            name=lang_info["name"],
        ))
    
    return sorted(supported_languages, key=lambda x: x.name)
