from uuid import uuid4

from django.db import models
# Create your models here.
from django.db.models import Model

from domain.books.entities import Book


class BookModel(Model):
    id = models.UUIDField(primary_key=True, default=uuid4, editable=False)
    author = models.TextField()
    title = models.TextField()
    publishing_date = models.DateTimeField()

    def to_entity(self) -> Book:
        return Book(id=self.id, author=self.author, title=self.title, publishing_date=self.publishing_date)
