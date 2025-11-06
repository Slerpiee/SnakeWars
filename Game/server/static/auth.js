const AUTH_SERVER_URL = "http://127.0.0.1:8000" //передавать через шаблонизатор
const LOGIN_ROUT = "/login"
const REGINSTER_ROUT= "/register"
function auth(data, isSignIn){
    if (isSignIn){
        Login(data)
        return
    }
    Register(data)
}

async function Login(data){
    headers = {
        'Content-Type': 'application/json',
    }
    body = JSON.stringify(data)

    try{
        const response = await fetch(AUTH_SERVER_URL+LOGIN_ROUT,{
            method: "POST",
            headers: headers,
            body: body
        })
        if (response.status != 200) {
            console.error("Server responded with status code: ", response.status, response.statusText)
            return
        }
        const res = await response.json()
        console.log("Result: ", res)
        sessionStorage.setItem("token", res["token_type"] + " " + res["access_token"])
        
    } catch (error){
        console.error("Error while fetching Auth Api: ", error)
    }
    
    

}

async function Register(data){
    headers = {
        'Content-Type': 'application/json',
    }
    body = JSON.stringify(data)
    try{
        const response = await fetch(AUTH_SERVER_URL+REGINSTER_ROUT,{
            method: "POST",
            headers: headers,
            body: body
        })
        if (response.status != 200) {
            console.error("Server responded with status code: ", response.status, response.statusText)
            return
        }
        const res = await response.json()
        console.log("Result: ", res)
        
    } catch (error){
        console.error("Error while fetching Auth Api: ", error)
    }
}