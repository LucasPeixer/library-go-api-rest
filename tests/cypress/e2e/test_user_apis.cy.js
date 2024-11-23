import 'cypress-plugin-api';
import '../utils/validarObjeto';

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
        });
    })

    it('get user', () => {
        cy.api({
            method: 'POST',
            url: 'http://localhost:8080/api/v1/users',
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {

            expect(response.status).to.equal(200)
            const listaUsuarios = response.body;
            const objetoEsperado = {
                id: 4,
                name: "Babau",
                cpf: "32514250056",
                phone: "(48)98484-5555",
                email: "empresa@gmail.com",
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
        const usuarioId = 4;

        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/users/activate/${usuarioId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "User has been successfully activated")

        });
    });
    it('Deactivate user', () => {
        const usuarioId = 4;

        cy.api({
            method: 'PUT',
            url: `http://localhost:8080/api/v1/users/deactivate/${usuarioId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.property('message', "User has been successfully deactivated")

        });
    });
    it('Delete user', () => {
        const usuarioId = 12;
        cy.api({
            method: 'DELETE',
            url: `http://localhost:8080/api/v1/users/delete/${usuarioId}`,
            headers: {
                Authorization: `Bearer ${authToken}`
            }
        }).then((response) => {
            expect(response.status).to.eq(200)
            expect(response.body).to.have.property('message', "User has been successfully deleted")
        });
    });

});
