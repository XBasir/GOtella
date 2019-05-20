package main
////////////////////////////////////////////////  IMPORTANDO LIBRERIAS   //////////////////////////////////////////
//importando librerias necesarias
import (
    "net"                              //servidores tcp udp y http
    "fmt"                              //Escribir en consolo
    "log"                              //Mostrar mensajes de errores
    "bufio"
    "os"
    "strings"
    "encoding/json"
)

/////////////////////////////////////////////////////// FUNCION ERRORES ////////////////////////////////////////
func checkError( err error) {
 if err != nil {
       log.Fatal(err)
  }
}

type Message struct{
    TypeMensaje string
    Mensaje string
    Nombre string
    Descripcion string

}


type ListMasterHost struct {
 direccion *net.UDPAddr
}
var ArrListMasterHost []ListMasterHost


type ListPeerHost struct {
 direccion *net.UDPAddr
}
var ArrListPeerHost []ListPeerHost


type ListPeerArchivo struct{
    nombre string
    descripcion string
}
var ArrListPeerArchivo []ListPeerArchivo

var connPeer *net.UDPConn

var connPeerActivo *net.UDPAddr

var connMasterActivo *net.UDPAddr



func MessageControl(conn *net.UDPConn) {
  message := make([]byte, 1024) //buffer de los mensajes
  for {
    //leer mensajes
    n, addr, err := conn.ReadFromUDP(message)
    checkError(err)

    var mensaje Message
    json.Unmarshal(message[:n], &mensaje)


    switch mensaje.TypeMensaje{
      case "conexionMaster":
          
          fmt.Print("\n")
          fmt.Println("Master Peer: ", addr)
          fmt.Println("Mensaje de Master Peer: ", mensaje.Mensaje)

          var master ListMasterHost
          master.direccion = addr
          ArrListMasterHost =  append(ArrListMasterHost, master)

          _, err = conn.WriteToUDP([]byte("Conexion Master Peer exitosa"), addr)
          checkError(err)

        break
      case "conexionPeer":

          fmt.Print("\n")
          fmt.Println("Peer: ", addr)
          fmt.Println("Mensaje de Peer: ", mensaje.Mensaje)

          var peer ListPeerHost
          peer.direccion = addr
          ArrListPeerHost =  append(ArrListPeerHost, peer)

          _, err = conn.WriteToUDP([]byte("Conexion Peer exitosa"), addr)
          checkError(err)
          
        break
      case "busquedaPeerArchivos":

          fmt.Print("\n")
          fmt.Println("Peer: ", addr)
          fmt.Println("Mensaje de Peer: ", mensaje.Mensaje)

          connPeerActivo = addr

          fmt.Print("\n")
          fmt.Println("Buscando archivo en Peers conectados")
          for _, item := range ArrListPeerHost {

            Pdir1:=item.direccion.String()
            Pdir2:=addr.String()

            if Pdir1 != Pdir2 {
              _, err = conn.WriteToUDP(message[:n], item.direccion)
              checkError(err)
            }
           
          }

          fmt.Println("Buscando archivo en Masters Peers conectados")
          for _, item := range ArrListMasterHost {
            
            dir1:=item.direccion.String()
            dir2:=addr.String()

            if dir1 != dir2 {
              _, err = conn.WriteToUDP(message[:n], item.direccion)
              checkError(err)
            }
          }

        break
      case "existe":

          fmt.Print("\n")
          fmt.Println("Peer: ", addr)
          fmt.Println("Mensaje de Peer: ",mensaje.Mensaje)
          fmt.Println("Nombre Archivo: ", mensaje.Nombre)
          fmt.Println("Descripcion: ", mensaje.Descripcion)

          _, err = conn.WriteToUDP(message[:n], connPeerActivo)
          checkError(err)

        break
      case "NoExiste":

          fmt.Print("\n")
          fmt.Println("Peer: ", addr)
          fmt.Println("Mensaje de Peer: ",mensaje.Mensaje)
          fmt.Println("Nombre Archivo: ", mensaje.Nombre)

          _, err = conn.WriteToUDP(message[:n], connPeerActivo)
          checkError(err)

        break
      default:
          fmt.Print("\n")
          fmt.Println("No se puede reconocer el mensaje entrante")
        break
    }
  }
}


func escucharPeer(conn *net.UDPConn) {
  for{
  message := make([]byte, 1024) //buffer de los mensajes
  n, addr, err := conn.ReadFromUDP(message)
  checkError(err)
  var mensaje Message
  json.Unmarshal(message[:n], &mensaje)

  if mensaje.TypeMensaje == "busquedaPeerArchivos"{
    var noExiste = "Nexiste"
    fmt.Print("\n")
    fmt.Println("Master Peer: ", addr)
    fmt.Println("Mensaje de Master Peer: ",mensaje.Mensaje)

    for _, item := range ArrListPeerArchivo {
      if (mensaje.Nombre == item.nombre){

        mensaje.TypeMensaje="existe"
        mensaje.Mensaje="Archivo encontrado"
        mensaje.Nombre=item.nombre
        mensaje.Descripcion=item.descripcion

        noExiste = "Existe"

        mensajeJson, _ := json.Marshal(mensaje)

        _, err = conn.WriteToUDP(mensajeJson,addr)
        checkError(err)
        break
      }
    }

    if noExiste != "Existe"  {

      mensaje.TypeMensaje="NoExiste"
      mensaje.Mensaje="Archivo no encontrado"

      mensajeJson, _ := json.Marshal(mensaje)

      _, err = conn.WriteToUDP(mensajeJson,addr)
      checkError(err)
    }
  }

 if mensaje.TypeMensaje == "existe" {
    fmt.Print("\n")
    fmt.Println("Master Peer: ", addr)
    fmt.Println("Mensaje de Master Peer: ",mensaje.Mensaje)
    fmt.Println("Nombre Archivo: ", mensaje.Nombre)
    fmt.Println("Descripcion: ", mensaje.Descripcion)
  }

   if mensaje.TypeMensaje == "NoExiste" {
    fmt.Print("\n")
    fmt.Println("Master Peer: ", addr)
    fmt.Println("Mensaje de Master Peer: ",mensaje.Mensaje)
    fmt.Println("Nombre Archivo: ", mensaje.Nombre)

  }

 }
}

func iniciarMaster() {

  reader := bufio.NewReader(os.Stdin)
  salir:=false
  for salir == false{
    fmt.Print("\n")
    fmt.Println("1-Crear master")
    fmt.Println("2-Crear master y unirse a na red p2p")
    fmt.Println("3-Salir")
    opt, _ := reader.ReadString('\n')
    switch opt{
      case "1\r\n":
          fmt.Print("\n")
          fmt.Println("Introdusca su direccion ip y puerto, ejemplo: 127.0.0.1:8080")
          fmt.Print("->")
          direccion, _ := reader.ReadString('\n')

          dir := strings.TrimRight(direccion, "\r\n")

          udpAddr, err := net.ResolveUDPAddr("udp4", dir)
          checkError(err)
          conn, err := net.ListenUDP("udp", udpAddr)
          checkError(err)
          defer conn.Close()


          fmt.Print("\n")
          fmt.Println("Servidor Peer Master: ", dir)

          MessageControl(conn)

        break
      case "2\r\n":
          fmt.Print("\n")
          fmt.Println("Introdusca su direccion ip y puerto, ejemplo: 127.0.0.1:8080")
          fmt.Print("->")
          direccion, _ := reader.ReadString('\n')

          dir := strings.TrimRight(direccion, "\r\n")

          udpAddr, err := net.ResolveUDPAddr("udp4", dir)
          checkError(err)
          conn, err := net.ListenUDP("udp", udpAddr)
          checkError(err)
          defer conn.Close()

          fmt.Print("\n")
          fmt.Println("Servidor Peer Master: ", dir)

          fmt.Print("\n")
          fmt.Println("Conectarse a, ejemplo: 127.0.0.1:8080")
          fmt.Print("->")
          direccion2, _ := reader.ReadString('\n')

          dirConexion := strings.TrimRight(direccion2, "\r\n")

          RemoteAddr, err := net.ResolveUDPAddr("udp", dirConexion)
          checkError(err)


          message := make([]byte, 1024) //buffer de los mensajes

          mensaje := Message{
            TypeMensaje: "conexionMaster",
            Mensaje: "Union Master Peer",
          }

          mensajeJson, _ := json.Marshal(mensaje)
          //escribir
          _, err = conn.WriteToUDP(mensajeJson,RemoteAddr)
          checkError(err)

          //leer mensajes
          n, addr, err := conn.ReadFromUDP(message)
          checkError(err)
          fmt.Print("\n")
          fmt.Println("Master Peer: ", addr)
          fmt.Println("Mensaje de Master Peer: ", string(message[:n]))

          var master ListMasterHost
          master.direccion = addr
          ArrListMasterHost = append(ArrListMasterHost, master)

          MessageControl(conn)

        break
      case "3\r\n":
          salir = true
        break
      default:
          fmt.Print("\n")
          fmt.Println("opcion invalida")
        break
    }
  }
}

func iniciarPeer(){
  reader := bufio.NewReader(os.Stdin)
  salir:=false
  for salir == false{
    fmt.Print("\n")
    fmt.Println("1-Crear cliente Peer y conectarse a la red P2P")
    fmt.Println("2-Crear archivo")
    fmt.Println("3-Lista de archivos")
    fmt.Println("4-Buscar archivos en la red P2P")
    fmt.Println("5-Salir")
    opt, _ := reader.ReadString('\n')

    switch opt{
      case "1\r\n":
          fmt.Print("\n")
          fmt.Println("Introdusca su direccion ip y puerto, ejemplo: 127.0.0.1:8080")
          fmt.Print("->")
          direccion, _ := reader.ReadString('\n')

          dir := strings.TrimRight(direccion, "\r\n")

          udpAddr, err := net.ResolveUDPAddr("udp4", dir)
          checkError(err)
          conn, err := net.ListenUDP("udp", udpAddr)
          checkError(err)
          defer conn.Close()

          connPeer = conn
          fmt.Print("\n")
          fmt.Println("Servidor Peer: ", dir)

          fmt.Print("\n")
          fmt.Println("Conectarse a, ejemplo: 127.0.0.1:8080")
          fmt.Print("->")
          direccion2, _ := reader.ReadString('\n')

          dirConexion := strings.TrimRight(direccion2, "\r\n")

          RemoteAddr, err := net.ResolveUDPAddr("udp", dirConexion)
          checkError(err)

          mensaje := Message{
            TypeMensaje: "conexionPeer",
            Mensaje: "Union Peer",
          }

          mensajeJson, _ := json.Marshal(mensaje)

          //escribir
          _, err = conn.WriteToUDP(mensajeJson,RemoteAddr)
          checkError(err)

          message := make([]byte, 1024) //buffer de los mensajes

          //leer mensajes
          n, addr, err := conn.ReadFromUDP(message)
          checkError(err)
          fmt.Print("\n")
          fmt.Println("Master Peer: ", addr)
          fmt.Println("Mensaje de Master Peer: ", string(message[:n]))

          connMasterActivo = addr

          go escucharPeer(conn)

        break
      case "2\r\n":
          fmt.Print("\n")
          fmt.Println("Introdusca el nombre del archivo")
          fmt.Print("->")
          nombre1, _ := reader.ReadString('\n')

          fmt.Print("\n")
          fmt.Println("Introdusca la descripcion del archivo")
          fmt.Print("->")
          descripcion1, _ := reader.ReadString('\n')

          nombre := strings.TrimRight(nombre1, "\r\n")
          direccion := strings.TrimRight(descripcion1, "\r\n")

          var archivo ListPeerArchivo
          archivo.nombre = nombre
          archivo.descripcion = direccion

          ArrListPeerArchivo =  append(ArrListPeerArchivo, archivo)

        break
      case "3\r\n":
        if(ArrListPeerArchivo != nil){
          i := 1
          for _, item := range ArrListPeerArchivo {
            fmt.Print("\n")
            fmt.Print("\n")
            fmt.Println("Archivo: ", i)
            fmt.Println("nombre : ", item.nombre)
            fmt.Println("descripcion : ", item.descripcion)
            fmt.Print("\n")
            i++
          }
        }else{
          fmt.Print("\n")
          fmt.Println("No hay archivos registrados!!!!")
        }
        break

      case "4\r\n":
          fmt.Print("\n")
          fmt.Println("Introdusca el nombre del archivo")
          fmt.Print("->")
          nombre1, _ := reader.ReadString('\n')

          nombre := strings.TrimRight(nombre1, "\r\n")

          mensaje := Message{
            TypeMensaje: "busquedaPeerArchivos",
            Mensaje: "Busqueda de archivo",
            Nombre: nombre,
          }

          mensajeJson, _ := json.Marshal(mensaje)

          //escribir
          _, err := connPeer.WriteToUDP(mensajeJson, connMasterActivo)
          checkError(err)

        break

      case "5\r\n":
          salir = true
        break
      default:
          fmt.Print("\n")
          fmt.Println("opcion invalida")
        break
    }
  }
}


func main() {
  salir:=false
  reader := bufio.NewReader(os.Stdin)

  for salir == false{
    fmt.Print("\n")
    fmt.Println("1-Master")
    fmt.Println("2-Peer")
    fmt.Println("3-Salir")
    opt, _ := reader.ReadString('\n')

    switch opt{
      case "1\r\n":
          iniciarMaster()
        break
      case "2\r\n":
          iniciarPeer()
        break
      case "3\r\n":
          salir = true
        break
      default:
          fmt.Print("\n")
          fmt.Println("opcion invalida")
        break
    }
  }

}