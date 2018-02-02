/* ==================================================================
CATALOG.GO
ppizz 2017 V0.1 photogrid
=====================================================================*/

package catalog

import (
    "fmt"
    "log"
    "strings"
    "io/ioutil"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

const VERSION = "Catalogue Photo V0.1"
const MAX = 300

type typCat struct {
	Id   int
	Name string
    DirName string
    Album string
	Title string
	Note int
    Label int
    Artist string
    Date string
    Modele string
	Orientation string
    Iso string
    Speed string
    Aperture string
    Focal string
    FileSize string
	Resolution string
    ExistPhoto bool
}
type typTabCat [MAX]typCat

var catDB *sql.DB
var	Tab typTabCat
var NbRecord int // photos dans DB
var NbPhoto int  // photos dans DIR

// ------------------------------------------------------------------------
//fonction renvoie la Version
func GetVersion () (vers string) {
    vers = VERSION
return
}

// ------------------------------------------------------------------------
//fonction rempli le tableau avec les noms de fichier jpg (sans les -th.jpg)
func InitTab () {
    for i := 0; i < MAX; i++ {
       Tab[i].Id = i
       Tab[i].Name = "-"
    }
}

// ------------------------------------------------------------------------
func GetTab() {
	fmt.Println("Catalog: \n", Tab[:NbRecord])
}

// ------------------------------------------------------------------------
func GetNbDirDB() (ret int) {
	var nbRecord int
    rows, err := catDB.Query("SELECT count(*) FROM Dir")
    if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
    rows.Next()
	err = rows.Scan(&nbRecord)
	if err != nil {
			log.Fatal(err)
			ret = -1
		} else {
	        ret = nbRecord
		}
	return	
}

// ------------------------------------------------------------------------
//Ouvre la base photo.db, qui doit se trouver dans Dir
func Init(Dir string) {
    var err error
    catDB, err = sql.Open("sqlite3", Dir + "photo.db")
    if err != nil {
		log.Fatal(err)
	 }
    InitTab()
}


// -------------------------------------------------------------
// Recherche dans la base les photos du catalogue 
// si OK met à jour tabPhoto[] et retourne le nombre de photo
// on laisse ouvert la base en sortie de fonction 
func GetPhotoDB(Dir string) (ret int) {
    var err error
    var rows *sql.Rows
    var maDate string

    i:=0
    rows, err = catDB.Query("select id, Name, DirName, Date, Album, Title, Note, Artist, Orientation  from photo WHERE dirname='"+Dir+"'")
    if err != nil {
        log.Fatal(err)
      }
    NbRecord = 0
	for rows.Next() {
        err = rows.Scan(&Tab[i].Id, &Tab[i].Name, &Tab[i].DirName, &maDate, &Tab[i].Album, &Tab[i].Title, &Tab[i].Note, &Tab[i].Artist, &Tab[i].Orientation)		
 		if err != nil {
			log.Fatal(err)
          }
       Tab[i].Date = maDate[8:10] + "/" + maDate[5:7] + "/" + maDate[:4]  
       Tab[i].ExistPhoto = false   
       i++
		}
    NbRecord = i   
    ret = NbRecord  

    return
}
 

// ------------------------------------------------------------------------
//fonction parcours le dossier pour les noms de fichier en .jpg
// ExistPhoto permet de verifier si le catalogue est synchronisé avec le dossier 
func GetPhotoDir(Dir string) {

	files, _ := ioutil.ReadDir(Dir)
	var sTab []string
    fname:= ""
    ext:= ""
    NbPhoto = 0
    for _, f := range files {	
	        sTab = strings.Split(f.Name(),".")
            fname = sTab[0]	
            ext = strings.ToUpper(sTab[1])
			if (ext == "JPG") {
                // on parcours le tableau à la recherche de la photo (les vignettes sont ignorées)
                if fname[len(fname)-3:] != "-th" {

                    for i := 0; i < NbRecord; i++ {
                        if Tab[i].Name == fname {
                            Tab[i].ExistPhoto = true
                            break
                         }
                    }
                   NbPhoto++
                }
		    }	
     }
}