from uuid import UUID
from abc import ABC, abstractmethod

from domain.books.entities import Book


class WriteBooksRepository(ABC):
    @abstractmethod
    def save_book(self, book: Book):
        ...


class ReadBooksRepository(ABC):
    @abstractmethod
    def get_book(self, _id: UUID) -> Book:
        ...


class BooksRepository(WriteBooksRepository, ReadBooksRepository, ABC):
    pass


class InMemoryRepository(BooksRepository):
    def save_book(self, book: Book):
        pass

    def get_book(self, _id: UUID):
        pass
