/**
 * Valida se um objeto esperado está presente em uma lista.
 * 
 * @param {Array} lista - A lista de objetos retornada pela API.
 * @param {Object} objetoEsperado - O objeto esperado que deve ser encontrado na lista.
 * @returns {boolean} - Retorna true se o objeto esperado for encontrado, caso contrário, false.
 */
export function validarObjetoNaLista(lista, objetoEsperado) {
    return lista.some((objeto) => {
        const valoresBasicosIguais =
            objeto.id === objetoEsperado.id &&
            objeto.title === objetoEsperado.title &&
            objeto.synopsis === objetoEsperado.synopsis &&
            objeto.amount === objetoEsperado.amount &&
            objeto.stock === objetoEsperado.stock;

        const autorIgual =
            objeto.author.id === objetoEsperado.author.id &&
            objeto.author.name === objetoEsperado.author.name;

        const generosIguais =
            objeto.genres.length === objetoEsperado.genres.length &&
            objeto.genres.every((genero) =>
                objetoEsperado.genres.some(
                    (g) => g.id === genero.id && g.name === genero.name
                )
            );
        
        return valoresBasicosIguais && autorIgual && generosIguais;
    });
}

export function validarObjetoNaListaStock(lista, objetoEsperado) {
    return lista.some((objeto) =>
        objeto.id === objetoEsperado.id &&
        objeto.status === objetoEsperado.status &&
        objeto.code === objetoEsperado.code &&
        objeto.book_id === objetoEsperado.book_id
    );
}