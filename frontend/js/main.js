const btnMenu = document.querySelector('#btn-menu')

//Função para mudar o menu hambúrguer
btnMenu.addEventListener('click', ()=>{
    
    const linhas = document.querySelectorAll('.hamburguer-line')
    const menu = document.querySelector('#menu')

    //Troca as classes tailwind ao clicar
    linhas[1].classList.toggle('w-3/4')
    linhas[2].classList.toggle('w-1/2')
    linhas[1].classList.toggle('w-full')
    linhas[2].classList.toggle('w-full')

    menu.classList.toggle('hidden')

})