from uuid import UUID

from books.models import BookModel
from domain.books.entities import Book
from domain.books.repositories import BooksRepository


class DjangoBooksRepository(BooksRepository):
    def get_book(self, _id: UUID) -> Book:
        book_model = BookModel.objects.get(id=_id)
        return book_model.to_entity()

    def save_book(self, book: Book) -> None:
        pass

