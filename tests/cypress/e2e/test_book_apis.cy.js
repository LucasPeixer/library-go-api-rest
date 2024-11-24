import 'cypress-plugin-api';
import {stock, genres, author, bookObject} from '../data/data_books'

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
                    20,21
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
            expect(response.body.stock).to.deep.equal(stock);
            expect(response.body.stock).to.deep.equal(author);
            expect(response.body.stock).to.deep.equal(genres);

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

            expect(response.status).to.equal(200)
            const listBooks = response.body;
            const findObject = validarObjetoNaLista(listBooks, bookObject);
            expect(findObject).to.be.true;

        })
    });
    it('Edit book', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/books/update/${userId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body:{
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
        const userId = Cypress.env('userId');
        cy.api({
            method: 'POST',
            url: `http://localhost:8080/api/v1/books/${userId}/stock/add`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body:{
                "code":22
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "Book stock added")

        });
    });
    it('Add book stock', () => {
        const userId = Cypress.env('userId');
        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/books/${userId}/stock/add`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body:{
                "code":22
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "Book stock added")
            const bookStockId = response.body.id;
            Cypress.env('bookStockId', bookStockId);

        });
    });


    it('get book stock', () => {
        cy.api({
            method: 'get',
            url: `http://localhost:8080/api/v1/books/`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            const listBooks = response.body;
            const findObject = validarObjetoNaLista(listBooks, bookObject);
            expect(findObject).to.be.true;

        })
    });

    it('Delete user', () => {
        const userId = Cypress.env('userId');
        cy.api({
            method: 'DELETE',
            url: `http://localhost:8080/api/v1/users/delete/${userId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            expect(response.body).to.have.property('message', "User has been successfully deleted")
        });
    });

});
