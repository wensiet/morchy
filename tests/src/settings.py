from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


def get_base_config(env_prefix: str) -> SettingsConfigDict:
    return SettingsConfigDict(
        env_file=".env",
        env_prefix=env_prefix,
        extra="ignore",
        env_file_encoding="utf-8",
    )


class ControlPlaneSettings(BaseSettings):
    url: str = Field(description="Controlplane API url")
    edge_url: str | None = Field(default=None, description="Edge proxy URL")
    seed_token: str | None = Field(default=None, description="Seed token for authentication")

    model_config = get_base_config("cp_")


class AgentSettings(BaseSettings):
    hosts: list[str] = Field(description="List of agent url's")
    ssh_key_file_path: str = Field(description="Path to ssh private key")

    model_config = get_base_config("agent_")


control_plane_settings = ControlPlaneSettings()
