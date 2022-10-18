from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional
from uuid import UUID


@dataclass
class Book:
    title: str
    author: str
    publishing_date: datetime
    id: Optional[UUID] = field(default=None)
