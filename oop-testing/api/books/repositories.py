from uuid import UUID

from books.models import BookModel
from domain.books.entities import Book
from domain.books.repositories import BooksRepository


class DjangoBooksRepository(BooksRepository):
    def get_book(self, _id: UUID) -> Book:
        try:
            book_model = BookModel.objects.get(id=_id)
        except BookModel.DoesNotExist:
            raise Exception("book not found")
        else:
            return book_model.to_entity()

    def save_book(self, book: Book) -> None:
        book_model = BookModel(author=book.author, publishing_date=book.publishing_date, title=book.title)
        book_model.save()

        book.id = book_model.id
