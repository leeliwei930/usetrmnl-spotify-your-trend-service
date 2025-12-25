from strands import Agent
from pydantic import BaseModel, Field
from strands.models.openai import OpenAIModel
import os
from strands.types.exceptions import StructuredOutputException

SYSTEM_PROMPT="""
You are an expert musicologist and linguist specializing in song meanings.
Your task is to provide concise, insightful summaries of the lyrics' meaning for a given song.

Instructions:
1. Analyze the song provided (Title, Artist, Album).
2. Distill the core message, themes, and emotional narrative.
3. Output strictly according to the requested schema (English and Mandarin summaries).
4. Adhere strictly to the word count limits (max 32 words per summary) to prevent token waste.
5. If the song is instrumental or has no lyrics, describe the musical mood instead.
"""

# Get LLM timeout from environment (default: 90 seconds)
LLM_TIMEOUT = int(os.getenv("LLM_TIMEOUT", "90"))


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




def invoke_agent(searchInput: LyricsSearchInput):
	# Sanitize inputs for cache key to ensure consistency
	clean_title = searchInput.title.strip().replace(' ', '')
	clean_artist = searchInput.artist.strip().replace(' ', '')
	clean_album = searchInput.album.strip().replace(' ', '')
	
	cache_key = f"lyrics_prompt_cache_{clean_title}_{clean_artist}_{clean_album}"

	model = OpenAIModel(
		model_id="openai/gpt-4.1:online",
		client_args={
			"base_url": "https://openrouter.ai/api/v1",
			"api_key": os.getenv("OPENROUTER_API_KEY"),
			"timeout": LLM_TIMEOUT,
		},
		params={
			"prompt_cache_key": cache_key,
			"reasoning_effort": "low",
		}
	)

	behind_the_lyrics_agent = Agent(
		model=model,
		system_prompt=SYSTEM_PROMPT,
	)
	response = behind_the_lyrics_agent(
		searchInput.model_dump_json(),
		structured_output_model=LyricSummary,
	)
	return response.structured_output
