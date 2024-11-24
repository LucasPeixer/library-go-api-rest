/**
 * Valida se um objeto esperado está presente em uma lista.
 * 
 * @param {Array} lista - A lista de objetos retornada pela API.
 * @param {Object} objetoEsperado - O objeto esperado que deve ser encontrado na lista.
 * @returns {boolean} - Retorna true se o objeto esperado for encontrado, caso contrário, false.
 */
export function validarObjetoNaLista(lista, objetoEsperado) {
    return lista.some((objeto) =>
        objeto.id === objetoEsperado.id &&
        objeto.title === objetoEsperado.title &&
        objeto.synopsis === objetoEsperado.synopsis &&
        objeto.amount === objetoEsperado.amount &&
        objeto.stock === objetoEsperado.stock &&
        objeto.author.id === objetoEsperado.author.id &&
        objeto.author.name === objetoEsperado.author.name &&
        objeto.genres.length === objetoEsperado.genres.length &&
        objeto.genres.every((genero, index) =>
            genero.id === objetoEsperado.genres[index].id &&
            genero.name === objetoEsperado.genres[index].name
        )
    );
}