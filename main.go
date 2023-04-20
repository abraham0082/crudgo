package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "abraham"
	password := "123456"
	basededatos := "sistema1"

	conexion, err := sql.Open(Driver, Usuario+":"+password+"@tcp(127.0.0.1)/"+basededatos)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

// esto es los template para hacer el crud de las platillas
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor Corriendo...")
	fmt.Println("Servidor corriendo2...")

	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id       int
	Nombre   string
	Correo   string
	Password string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hola Sublimamani")
	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
		//fmt.Println("error al conectarse")
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, password, correo string

		//err = registros.Scan(&id, &nombre, &correo)
		err = registros.Scan(&id, &nombre, &password, &correo)
		if err != nil {
			http.Error(w, "Error en la consulta", http.StatusInternalServerError)
			return
			//fmt.Println("error con los usuarios")
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Password = password
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
	//plantillas.ExecuteTemplate(w, "crear", nil)
	//fmt.Fprintf(w, "estoy en crear")
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	//fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	BorrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	BorrarRegistro.Exec(idEmpleado)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Editar(w http.ResponseWriter, r *http.Request) {

	idEmpleado := r.URL.Query().Get("id")
	//fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)
	if err != nil {
		panic(err.Error())
		//fmt.Println("error al conectarse")
	}

	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, password, correo string

		//err = registros.Scan(&id, &nombre, &correo)
		err = registro.Scan(&id, &nombre, &password, &correo)
		if err != nil {
			http.Error(w, "Error en la consulta", http.StatusInternalServerError)
			return
			//fmt.Println("error con los usuarios")
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Password = password
		empleado.Correo = correo

	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		password := r.FormValue("password")
		correo := r.FormValue("correo")
		fmt.Println(id)
		fmt.Println(nombre)
		fmt.Println(password)
		fmt.Println(correo)

		conexionEstablecida := conexionBD()
		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, password=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		actualizarRegistro.Exec(nombre, password, correo, id)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
