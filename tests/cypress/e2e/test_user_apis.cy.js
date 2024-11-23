import 'cypress-plugin-api';
import {validarObjetoNaLista

} from '../utils/validarObjeto';

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

    it('create user', () => {
        cy.api({
            method: 'POST',
            url: 'http://localhost:8080/api/v1/register',
            headers: {
                Authorization: `Bearer ${authToken}`,
            },
            body: {
                "name": "usuário",
                "cpf": "64801920012",
                "phone": "(48)98484-6666",
                "Email": "empresa22@gmail.com",
                "password": "123",
                "role_id": 2
            }
        }).then((response) => {

            expect(response.status).to.equal(201)
            expect(response.body).to.have.property('message', "User registered successfully")
            const userId = response.body.user_id;  // Verifique o nome exato da chave
            Cypress.env('userId', userId);
        });
    })

    it('get user', () => {
        const userId = Cypress.env('userId');
        cy.log(userId)
        cy.api({
            method: 'get',
            url: 'http://localhost:8080/api/v1/users',
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            const listaUsuarios = response.body;
            const objetoEsperado = {
                id: userId,
                name: "usuário",
                cpf: "64801920012",
                phone: "(48)98484-6666",
                email: "empresa22@gmail.com",
                account_role: {
                    id: 2,
                    name: "user"
                },
                is_active: true,
            };
            const objetoEncontrado = validarObjetoNaLista(listaUsuarios, objetoEsperado);
            expect(objetoEncontrado).to.be.true;

        })
    });
    it('Activate user', () => {
        const userId = Cypress.env('userId');
        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/users/activate/${userId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "User has been successfully activated")

        });
    });
    it('Deactivate user', () => {
        const userId = Cypress.env('userId');
        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/users/deactivate/${userId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "User has been successfully deactivated")

        });
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
