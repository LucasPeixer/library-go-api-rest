/**
 * Função para verificar se um objeto esperado está presente na lista de usuários do sistema.
 * 
 * @param {Array} lista - A lista de objetos retornada pela API.
 * @param {Object} objetoEsperado - O objeto a ser validado na lista.
 * @returns {Boolean} - Retorna `true` se o objeto estiver na lista, caso contrário `false`.
 */

export function validarObjetoNaLista(lista, objetoEsperado) {
    return lista.some((usuario) => 
      usuario.id === objetoEsperado.id &&
      usuario.name === objetoEsperado.name &&
      usuario.cpf === objetoEsperado.cpf &&
      usuario.phone === objetoEsperado.phone &&
      usuario.email === objetoEsperado.email &&
      usuario.account_role.id === objetoEsperado.account_role.id &&
      usuario.account_role.name === objetoEsperado.account_role.name &&
      usuario.is_active === objetoEsperado.is_active
    );
  }