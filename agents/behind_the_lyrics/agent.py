from strands import Agent
from pydantic import BaseModel, Field
from strands.models.openai import OpenAIModel
import os

SYSTEM_PROMPT="""
You're are a music expert agent, based on the given song name, artist and album detail.

"""

# Get LLM timeout from environment (default: 90 seconds)
LLM_TIMEOUT = int(os.getenv("LLM_TIMEOUT", "90"))

model = OpenAIModel(
    model_id="openai/gpt-4.1:online",
	client_args={
		"base_url": "https://openrouter.ai/api/v1",
        "api_key": os.getenv("OPENROUTER_API_KEY"),
        "timeout": LLM_TIMEOUT,
    }
)

class LyricsSearchInput(BaseModel):
	title: str = Field(description="The song's title")
	artist: str = Field(description="The song's artist")
	album: str = Field(description="The song's album")
	cover: str = Field(description="The song's cover url")

class LyricSummary(BaseModel):
	title: str = Field(description="The song's title")
	artist: str = Field(description="The song's artist")
	album: str = Field(description="The song's album")
	cover: str = Field(description="The song's cover url")
	summary_en: str = Field(description="A summarised version of the song's lyrics meaning, in 32 words")
	summary_zh: str = Field(description="A summarised version of the song's lyrics meaning in mandarin, in 32 words")


behind_the_lyrics_agent = Agent(
    model=model,
    system_prompt=SYSTEM_PROMPT,
)

def invoke_agent(searchInput: LyricsSearchInput):
	response = behind_the_lyrics_agent(
		searchInput.model_dump_json(),
		structured_output_model=LyricSummary,
	)
	return response.structured_output
