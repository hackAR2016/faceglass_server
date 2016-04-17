package main

import (
    "fmt"
    "log"
    "io"
    "net/http"
    "strconv"
    "encoding/json"
    "github.com/gorilla/mux"
    "os"
)

type User struct {
    ID        int
    Name      string
    Text      string
    Status    string
    Email     string
}

type Users []User

var users Users = Users{
    User{ID: 0, Name: "dario", Text: "programmer, backend", Status: "DnD", Email: ""},
    User{ID: 1, Name: "alexander", Text: "programmer, frontend", Status: "Come talk to me", Email: ""},
    User{ID: 2, Name: "yan_wo", Text: "programmer, frontend", Status: "DnD", Email: ""},
    User{ID: 3, Name: "jenny_li", Text: "designer", Status: "BrB", Email: ""},
}

func main() {
    initModel() 
    router := mux.NewRouter().StrictSlash(true)
    
    router.PathPrefix("/asset/").Handler( http.StripPrefix("/asset/", http.FileServer(http.Dir("./asset/"))) )
  
    router.HandleFunc("/label", getLabelHandler)
    router.HandleFunc("/users", userIndex).Methods("GET")
    router.HandleFunc("/users", addUser).Methods("POST")
    router.HandleFunc("/users/{userId}", userShow).Methods("GET")
    router.HandleFunc("/users/{userId}", changeUser).Methods("POST")
 
    os.MkdirAll("./asset/users", 0777)
    
    fmt.Println("Serving content at :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func getLabelHandler(w http.ResponseWriter, r *http.Request) {
    if err := json.NewEncoder(w).Encode(users[getLabel(0)]); err != nil {
        panic(err)
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

// upload logic
func changeUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    r.ParseMultipartForm(32 << 20)
    fmt.Println(r.Form.Encode());
    file, handler, err := r.FormFile("image")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    
    
    fmt.Fprintf(w, "%v", handler.Header)
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["userId"])
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    var foundUser User
    
    foundUser.ID = -1
    var foundUserIndex = -1
    
    for index,user := range users {
        if user.ID == userID {
            foundUser = user
            foundUserIndex = index;
            break
        }
    } 
    
    if foundUser.ID == -1 {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    users[foundUserIndex].Name = r.FormValue("user_name");
    users[foundUserIndex].Text = r.FormValue("user_comment");
    users[foundUserIndex].Email = r.FormValue("user_email");
    users[foundUserIndex].Status = "Available";
    

    userIDStr := strconv.Itoa(userID)
    
    os.MkdirAll("./asset/users/" + userIDStr, 0777);
    filename := "./asset/users/" + userIDStr + "/" + handler.Filename;
    f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    io.Copy(f, file)
    updateModel(userID, filename);
}

// upload logic
func addUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    r.ParseMultipartForm(32 << 20)
    fmt.Println(r.Form.Encode());
    file, handler, err := r.FormFile("image")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    
    
    fmt.Fprintf(w, "%v", handler.Header)

    userID := len(users)
    
    var newUser User;
    newUser.ID = userID;
    newUser.Name = r.FormValue("user_name");
    newUser.Text = r.FormValue("user_comment");
    newUser.Email = r.FormValue("user_email");
    newUser.Status = "Available";
    
    users = append(users, newUser)   

    userIDStr := strconv.Itoa(userID)
    
    os.MkdirAll("./asset/users/" + userIDStr, 0777);
    filename := "./asset/users/" + userIDStr + "/" + handler.Filename;
    f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    io.Copy(f, file)
    updateModel(userID, filename);
}

func userIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(users); err != nil {
        panic(err)
    }
}

func userShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["userId"])
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    var foundUser User
    
    foundUser.ID = -1
    
    for _,user := range users {
        if user.ID == userID {
            foundUser = user
            break
        }
    } 
    
    if foundUser.ID == -1 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(foundUser); err != nil {
        panic(err)
    }
}
