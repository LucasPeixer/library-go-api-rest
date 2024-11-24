import 'cypress-plugin-api';
import { genres, author} from '../data/data_books'
import { validarObjetoNaLista, validarObjetoNaListaStock } from '../utils/validarBookObject'

let authToken;
describe('API tests', () => {
    before(() => {
        cy.api({
            method: 'POST',
            url: 'http://localhost:8080/api/v1/login',
            body: {
                "email": "locodoparanaue46@gmail.com",
                "password": "123"
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            authToken = response.body.replace(/^"|"$/g, '');
            cy.log(`Auth Token: ${authToken}`);

        });
    });

    it('create book', () => {
        cy.api({
            method: 'POST',
            url: 'http://localhost:8080/api/v1/books/create',
            headers: {
                Authorization: `Bearer ${authToken}`,
            },
            body: {
                "title": "A casa amarela",
                "synopsis": "Uma casa que um dia foi amarela",
                "book_codes": [
                    20, 21
                ],
                "author_id": 2,
                "genre_ids": [
                    1, 2, 3, 4, 5
                ]
            }
        }).then((response) => {

            expect(response.status).to.equal(201)
            expect(response.body).to.have.property('title', "A casa amarela")
            expect(response.body).to.have.property('synopsis', "Uma casa que um dia foi amarela")
            expect(response.body).to.have.property('amount', 2)
            expect(response.body.author).to.deep.equal(author);
            expect(response.body.genres).to.deep.equal(genres);

            const bookId = response.body.id;
            Cypress.env('bookId', bookId);


        });
    })

    it('get book', () => {
        cy.api({
            method: 'get',
            url: `http://localhost:8080/api/v1/books/`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            const bookId = Cypress.env('bookId');
            const bookObject = {
                id: bookId,
                title: "A casa amarela",
                synopsis: "Uma casa que um dia foi amarela",
                amount: 2,
                stock: null,
                author: {
                    id: 2,
                    name: "José"
                },
                genres: [
                    { id: 3, name: "Mistério" },
                    { id: 4, name: "Fantasia" },
                    { id: 2, name: "Aventura" },
                    { id: 5, name: "Ficção Científica" },
                    { id: 1, name: "Romance" }
                ]
            };

            const listBooks = response.body;
            const findObject = validarObjetoNaLista(listBooks, bookObject);
            cy.log(listBooks)
            cy.log(bookId)
            cy.log(findObject)
            expect(response.status).to.equal(200)
            expect(findObject).to.be.true;

        })
    });
    it('Edit book', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/books/update/${bookId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body: {
                "title": "Baia amarela",
                "synopsis": "Um vale amarelo",
                "amount": 22,
                "author_id": 3
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "Book updated successfully")

        });
    });
    it('Add book stock', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'POST',
            url: `http://localhost:8080/api/v1/books/${bookId}/stock/add`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body: {
                "code": 22
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "Book stock added")
            const bookStockId = response.body.book_stock_id;
            Cypress.env('bookStockId', bookStockId);
        });
    });

    it('get book stock', () => {
        const bookId = Cypress.env('bookId');

        cy.api({
            method: 'get',
            url: `http://localhost:8080/api/v1/books/${bookId}/stock`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            const listBookStock = response.body;
            const stockObject = listBookStock[2];
            const findObject = validarObjetoNaListaStock(listBookStock, stockObject);
            cy.log(listBookStock)
            cy.log(findObject)
            expect(findObject).to.be.true;

        })
    });

    it('Delete book stock', () => {
        const bookId = Cypress.env('bookId');
        const bookStockId = Cypress.env('bookStockId');
        cy.api({
            method: 'DELETE',
            url: `http://localhost:8080/api/v1/books/${bookId}/stock/remove/${bookStockId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            expect(response.body).to.have.property('message', "Book stock removed")
        });
    });

    it('Delete book', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'DELETE',
            url: `http://localhost:8080/api/v1/books/delete/${bookId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            expect(response.body).to.have.property('message', "Book deleted successfully")
        });
    });

});
