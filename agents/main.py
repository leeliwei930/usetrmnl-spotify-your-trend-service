from pathlib import Path
from dotenv import load_dotenv
from fastapi import FastAPI
from behind_the_lyrics.agent import invoke_agent, LyricsSearchInput
from fastapi import Body
from typing import Annotated

# Load environment variables from .env file
env_path = Path(__file__).parent / ".env"
load_dotenv(dotenv_path=env_path)

app = FastAPI()

@app.get("/health")
async def health_check():
    return {"message": "ok"}

@app.post("/usetrmnl/agent/behind-the-lyrics")
async def agent_perform_lyrics_search(searchInput: Annotated[LyricsSearchInput, Body(embed=True)]):
    return invoke_agent(searchInput)
