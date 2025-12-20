from pathlib import Path
from dotenv import load_dotenv
from fastapi import FastAPI
from behind_the_lyrics.agent import invoke_agent, LyricsSearchInput
from fastapi import Body
from typing import Annotated
import os
import asyncio

# Load environment variables from .env file
env_path = Path(__file__).parent / ".env"
load_dotenv(dotenv_path=env_path)

# Get timeout from environment (default: 120 seconds)
REQUEST_TIMEOUT = int(os.getenv("REQUEST_TIMEOUT", "120"))

app = FastAPI()

@app.get("/health")
async def health_check():
    return {"message": "ok"}

@app.post("/usetrmnl/agent/behind-the-lyrics")
async def agent_perform_lyrics_search(searchInput: Annotated[LyricsSearchInput, Body(embed=True)]):
    # Run the synchronous agent call in a thread pool with timeout
    try:
        result = await asyncio.wait_for(
            asyncio.to_thread(invoke_agent, searchInput),
            timeout=REQUEST_TIMEOUT
        )
        return result
    except asyncio.TimeoutError:
        return {
            "error": "Request timeout - LLM agent took too long to respond",
            "timeout": REQUEST_TIMEOUT
        }
