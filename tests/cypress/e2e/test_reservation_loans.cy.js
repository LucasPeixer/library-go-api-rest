import 'cypress-plugin-api';
import { genres, author} from '../data/data_books'
import { validarObjetoNaListaReservas } from '../utils/validarBookObject'

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
                "title": "Casa amarela 1",
                "synopsis": "Uma casa que um dia foi amarela 1",
                "author_id": 2,
                "genre_ids": [
                    1, 2, 3, 4, 5
                ]
            }
        }).then((response) => {

            expect(response.status).to.equal(201)
            expect(response.body).to.have.property('title', "Casa amarela 1")
            expect(response.body).to.have.property('synopsis', "Uma casa que um dia foi amarela 1")
            expect(response.body).to.have.property('amount', 0)
            expect(response.body.author).to.deep.equal(author);
            expect(response.body.genres).to.deep.equal(genres);

            const bookId = response.body.id;
            Cypress.env('bookId', bookId);


        });
    })

    it('Add book stock', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'POST',
            url: `http://localhost:8080/api/v1/books/${bookId}/stock/add`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body: {
                "code": 100
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "Book stock added")
            const bookStockId = response.body.book_stock_id;
            Cypress.env('bookStockId', bookStockId);
        });
    });

    it('Create reservation', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'POST',
            url: `http://localhost:8080/api/v1/reservations/create`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body: {
                "book_id": bookId,
                "borrowed_days": 30
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            const objectResponse = response.body
            Cypress.env('responseReservation', objectResponse);
        });
    });
    
    it('get reservation', () => {
        const responseReservation = Cypress.env('responseReservation');

        cy.api({
            method: 'get',
            url: `http://localhost:8080/api/v1/reservations`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            const listReservation = response.body
            const findObject = validarObjetoNaListaReservas(listReservation, responseReservation);
            cy.log(responseReservation)
            cy.log(listReservation)
            cy.log(findObject)
            //expect(findObject).to.be.true;

        })
    });

    it('Create loan', () => {
        const bookId = Cypress.env('bookId');
        cy.api({
            method: 'POST',
            url: `http://localhost:8080/api/v1/loans/create`,
            headers: {
                Authorization: `Bearer ${authToken}`
            },
            body: {
                "book_id": bookId,
                "borrowed_days": 90
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            const objectResponse = response.body
            Cypress.env('responseLoan', objectResponse);
        });
    });

});
