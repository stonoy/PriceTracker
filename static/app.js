// the base url
const url = "/v1/"

// store the products
let products = []

// get the elements

// login elements
const loginDiv = document.getElementById("loginDiv")
const loginBtn = document.getElementById("loginBtn")
const emailLogin = document.getElementById("emailLogin")
const passwordLogin = document.getElementById("passwordLogin")
const errorMsgLogin = document.getElementById("errorMsgLogin")
const toggleLogin = document.getElementById("toggleLogin")

// register elements
const registerDiv = document.getElementById("registerDiv")
const registerBtn = document.getElementById("registerBtn")
const nameReg = document.getElementById("nameReg")
const emailReg = document.getElementById("emailReg")
const passwordReg = document.getElementById("passwordReg")
const errorMsgRegister = document.getElementById("errorMsgRegister")
const toggleReg = document.getElementById("toggleReg")

// other elements
const todoListContainer = document.getElementById("todoListContainer")
const closeModal = document.getElementById("closeModal")
const plus = document.getElementById("plus")
const modal = document.getElementById("modal")
const productName = document.getElementById("productName")
const productUrl = document.getElementById("productUrl")
const addProductBtn = document.getElementById("addProductBtn")
const logoutBtn = document.getElementById("logout")
const username = document.getElementById("username")

// const dropDownBtns = document.querySelectorAll("dropDownBtns")

// provide html for each product
function productHtml({name, url, id, created_at, updated_at, base_price, current_price, priority}) {
    return `<div class="relative ${priority ? 'border-red-500' : 'border-gray-300'} border-2 rounded">
        <div class="absolute p-2 top-0 right-0">
            <div class="dropdown inline-block relative">
                <button  class="bg-gray-200 text-gray-700 font-semibold py-2 px-4 rounded inline-flex items-center">
                    <i id="dropDown" class="text-md fas fa-caret-down"></i>
                </button>
                <ul class="dropdown-menu bg-slate-500 absolute hidden text-gray-700 pt-1">
                    <li id="edit" data-edit="${id}" class="bg-gray-200 hover:bg-gray-400 py-2 px-4 block whitespace-no-wrap">${priority ? "Unfollow" : "Follow"}</li>
                    <li id="delete" data-delete="${id}" class="bg-gray-200 hover:bg-gray-400 py-2 px-4 block whitespace-no-wrap">Delete</li>
                </ul>
            </div>
        </div>
        <a href="${url}" target="_blank" rel="noopener noreferrer" class="block bg-gray-200 p-4 rounded-md mb-2">
            <p class="text-lg font-semibold mb-1">${name}</p>
            <p class="text-sm font-bold ${current_price <= base_price ? "text-green-600" : "text-red-600"} ">${calculatePriceChange(base_price,current_price).toFixed(2)}%</p>
            <p class="text-sm text-gray-600 mb-1">Created : ${calculateTime(created_at, updated_at)} ago </p>
            <p class="text-sm font-medium text-gray-600">Base Price : ${base_price}</p>
            <p class="text-sm font-medium text-gray-600">Current Price : ${current_price}</p>
        </a>
    </div>`;
}


// EVENT LISTENERS...

// toggle login / register div
toggleLogin.addEventListener("click", toggleLR)
toggleReg.addEventListener("click", toggleLR)

// get add product modal
plus.addEventListener("click", ()=>{
    modal.classList.remove("hidden")
})

// cancel button
closeModal.addEventListener("click", ()=>{
    modal.classList.add("hidden")
})

// handle the down arrow key
todoListContainer.addEventListener("click", async (e)=> {
    if(e.target && e.target.id === "dropDown"){
        
        e.target.parentElement.nextElementSibling.classList.toggle("hidden")
    }

    if(e.target && e.target.id === "edit"){
        // console.log("edit", e.target.dataset.edit)
        await updateProduct(e.target.dataset.edit)
        
    }

    if(e.target && e.target.id === "delete"){
        // console.log("delete", e.target.dataset.delete)
        await deleteProduct(e.target.dataset.delete)
    }
})

// get data from register form and send it to server and toggle back to login
registerBtn.addEventListener("click", async (e)=>{
    e.preventDefault()

    // set error msg blank
    errorMsgRegister.innerText = ""

    let name = nameReg.value
    let email = emailReg.value
    let password = passwordReg.value
    
    if (!email || !password || !name){return}

    const user = await registerUser({name, email, password})
    
    emailReg.value = ""
    passwordReg.value = ""

    // console.log(user)

    if (!user){return}

    // toggle back to login div
    toggleLR()
})


// get data from login form and send it to server and set localstorage and fetch products
loginBtn.addEventListener("click", async (e)=>{
    e.preventDefault()

    // set error msg blank
    errorMsgLogin.innerText = ""

    let email = emailLogin.value
    let password = passwordLogin.value
    
    if (!email || !password){return}

    const user = await loginUser({email, password})
    
    emailLogin.value = ""
    passwordLogin.value = ""

    // console.log(user)

    if (!user){return}

    // set token to local storage
    localStorage.setItem("user", JSON.stringify(user))

    // eable the logout button
    logoutBtn.disabled = false

    // disable the register toggler button
    toggleLogin.disabled = true

    // show the username
    username.innerText = `Welcome ${user.name}`

     products = await getProducts(user.token)
    // console.log(products)

    renderProducts(products)
})

// post request to server to add a product
addProductBtn.addEventListener("click", async ()=>{
    if (!productUrl.value || !productName.value){return}

    const token = JSON.parse(localStorage.getItem("user"))?.token
    if(!token){
        alert("login to add products")
        return
    }

    const addInfo = await addProducts({name:productName.value, url:productUrl.value}, token)

    // console.log(addInfo)

    // set add products feilds empty
    productName.value = ""
    productUrl.value = ""

    // hide the modal and fetch new products
    

     products = await getProducts(token)
    // console.log(products)

    renderProducts(products)

    modal.classList.add("hidden")

})

logoutBtn.addEventListener("click", () => {
    // delete the token from local storage
    localStorage.clear()

    // refresh the page
    location.reload()

})

// register login
async function registerUser(registerObj){
    try {
        const resp = await fetch(`${url}register`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
              },
            body: JSON.stringify(registerObj),
        })
        // console.log(resp)
        if (resp.status == "404"){return}

        const data = await resp.json()

        if ("error" in data){
            let errMsg = data.error.split(":")[0]
            errorMsgRegister.innerText = errMsg
            // console.log(data)
            return
            
        }
        return data
    } catch (error) {
        
        console.log(error)
    }
}

// login logic
async function loginUser(loginObj){
    try {
        const resp = await fetch(`${url}login`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
              },
            body: JSON.stringify(loginObj),
        })
        // console.log(resp)
        if (resp.status == "404"){return}

        const data = await resp.json()

        if ("error" in data){
            let errMsg = data.error.split(":")[0]
            errorMsgLogin.innerText = errMsg
            // console.log(data)
            return
            
        }
        return data
    } catch (error) {
        
        console.log(error)
    }
}


// get request to server to fetch products
async function getProducts(token){
    try {
        const resp = await fetch(`${url}userproducts`, {
            method: "GET",
            headers : {
                "Authorization" : `Bearer ${token}`
            }
        })

        if (resp.status == "404"){return}

        const data = await resp.json()

        if ("error" in data){
            let errMsg = data.error.split(":")[0]
            alert(errMsg)
            // console.log(data)
            return
            
        }
        return data

    } catch (error) {
        console.log(error)
    }
}

// add products logic
async function addProducts(productObj, token){
    try {
        const resp = await fetch(`${url}createproducts`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(productObj)
        }) 
 
        if (resp.status == "404"){return}
 
        const data = await resp.json()

        if ("error" in data){
            alert(data.error.split(":")[0])
            console.log(data)
            return
        }
        return data
     } catch (error) {
         console.log(error)
     }
}

async function updateProduct(id){
    // check user is present (token)
    const token = JSON.parse(localStorage.getItem("user"))?.token
    if (!token){return}

    // send put request
    try {
      const resp = await fetch(`${url}updatepriority/${id}`, {
        method: "PUT",
        headers : {
            "Authorization": `Bearer ${token}`,
        }
      })

      if (resp.status == "404"){return}
 
        const data = await resp.json()

        if ("error" in data){
            alert(data.error.split(":")[0])
            console.log(data)
            return
        }
    } catch (error) {
        console.log(error)
    }

    // get the products
    products = await getProducts(token)
    // console.log(products)

    // render new products
    renderProducts(products)
}

async function deleteProduct(id){
     // check user is present (token)
     const token = JSON.parse(localStorage.getItem("user"))?.token
     if (!token){return}
 
     // send delete request
     try {
       const resp = await fetch(`${url}deleteproduct/${id}`, {
         method: "DELETE",
         headers : {
             "Authorization": `Bearer ${token}`,
         }
       })
 
       if (resp.status == "404"){return}
  
         const data = await resp.json()
 
         if ("error" in data){
             alert(data.error.split(":")[0])
             console.log(data)
             return
         }
     } catch (error) {
         console.log(error)
     }
 
     // get the products
     products = await getProducts(token)
    //  console.log(products)
 
     // render new products
     renderProducts(products)
}

// render products in container
function renderProducts({products}){
    let productsHtml = ""

    products.forEach((product) => {
        productsHtml += productHtml(product)
    })

    todoListContainer.innerHTML = productsHtml
}

// re render when page loads
async function render(){
    const {token,name} = JSON.parse(localStorage.getItem("user"))
    // no user then return
    if (!token){return}

    // eable the logout button
    logoutBtn.disabled = false

    // disable the register toggler button
    toggleLogin.disabled = true

    // show the username
    username.innerText = `Welcome ${name}`

    // any user toggle login div and show products
    toggleLR()
     products = await getProducts(token)
    // console.log(products)

    renderProducts(products)
}

// comeback main page with get request to server to fetch products

render()

// Utils...

// toggle login/register div
function toggleLR(){
    loginDiv.classList.toggle("hidden")
    registerDiv.classList.toggle("hidden")
}

// calculate price change
function calculatePriceChange(base, current){
    return (current - base)*100/base
}

// calculate time difference
function calculateTime(start, now){

const startDate = new Date(start);
const endDate = new Date(now);

const timeDifference = Math.abs(endDate - startDate);

// Convert time difference to days, hours, minutes, and seconds
const days = Math.floor(timeDifference / (1000 * 60 * 60 * 24));
const hours = Math.floor((timeDifference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
const minutes = Math.floor((timeDifference % (1000 * 60 * 60)) / (1000 * 60));
const seconds = Math.floor((timeDifference % (1000 * 60)) / 1000);

return `${days} days, ${hours} hours, ${minutes} minutes, ${seconds} seconds`

}