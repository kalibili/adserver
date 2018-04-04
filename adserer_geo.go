package main


import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "io/ioutil"
    "html/template"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
    "strconv"
)

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en">
<head></head>
<body>
{{ range $value := . }}
   <img src ="{{ $value }}" alt="state.pic" height="256" width="453 ">
{{ end }}
</body>`

/*type client struct {
    Id int
    client_no string
    name string
}

type ad_campaign struct {
    Id int
    client *client
    name string
    targeting_type_pos bool
}
*/

type state struct {
    state_code string
}

type ad_pos_neg struct {
	ad_campaign_id int
    state_ids []string
    targeting_type_pos bool
}

type images struct {
    image_location string
    ad_campaign_id int
}

/*client_map := make(map[int]struct)
ad_campaign_map := make(map[int]struct)*/
var state_map = make(map[string]state)
var ad_pos_neg_map = make(map[string]ad_pos_neg)
var images_map = make(map[string]images)

func getRegion(w http.ResponseWriter, r *http.Request) {
	var data interface{}
    hostname := r.RemoteAddr
    ip_addr := strings.Split(hostname,":")[0]
    fmt.Println("Request from:",ip_addr)
    if ip_addr == "127.0.0.1" {
        fmt.Println("Geo location not enabled when server is running locally")
    }
    url := fmt.Sprintf("http://ip-api.com/json/%s",ip_addr)
    //url := fmt.Sprintf("http://ip-api.com/json/1.186.78.10")
    res, err := http.Get(url)
    if err != nil {
        panic(err.Error())
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err.Error())
    }

    json.Unmarshal(body, &data)
    st, _ := data.(map[string]interface{})
    fmt.Println("Current state:",st["region"])
    var pos_ad []string

    for a, b := range state_map {                                                              //a:id of state b:state_map
        if b.state_code==st["region"] {
            for _,d := range ad_pos_neg_map {                                                  //c:id of ad_pos_neg_map d:ad_pos_neg_map
                for _,e :=range d.state_ids {                                                  //e:values in state_ids i.e id of states
                    if a == e && d.targeting_type_pos == true {                                //to check whether the current st id is present in ad_pos_neg state_id's
                        pos_ad = append(pos_ad,strconv.Itoa(d.ad_campaign_id))
                    }
                }

            }
        }
    }
    var neg_ad []string
    var State_present bool

    for a, b := range state_map {                                                              //a:id of state b:state_map
        if b.state_code==st["region"] {
            for _,d := range ad_pos_neg_map {                                                  //c:id of ad_pos_neg_map d:ad_pos_neg_map
                State_present = false
                if d.targeting_type_pos == false {
                    for _,e :=range d.state_ids {                                              //e:values in state_ids i.e id of states
                        if a == e {                                                            //to check whether the current st id is present in ad_pos_neg state_id's
                            State_present = true
                            break
                        }
                        }

                    if State_present == false {
                        neg_ad = append(neg_ad,strconv.Itoa(d.ad_campaign_id))
                    }
                }
                }

            }
        }

    for _,j :=range neg_ad {
        pos_ad = append(pos_ad,j)
    }

    var img_string []string

    for _,g :=range images_map {
        for _,h :=range pos_ad {
            if strconv.Itoa(g.ad_campaign_id) == h {
                fmt.Println("Images  ",g.image_location)
                img_string =append(img_string,g.image_location)
            }
        }
    }




    t := template.New("t")
    t, err1 := t.Parse(ImageTemplate)
    if err1 != nil {
        panic(err)
    }
    err = t.Execute(w, img_string)
    if err != nil {
        panic(err)
    }

}


func main() {
	fmt.Println("Loading...")
    dbinfo := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable","kalibili", "adserer2","kalibili" )
    db, err := sql.Open("postgres", dbinfo)
    if err != nil {
        log.Fatal(err)
    }
    /*ad_campaign_sql,err_ad_campaign := db.Query("select id from ad_campaign")
    if err_ad_campaign != nil {
        log.Fatal(err_ad_campaign)
    }*/
    state_sql,err_state := db.Query("select id,state_code from state")
    if err_state != nil {
        log.Fatal(err_state)
    }
    ad_pos_neg_sql,err_ad_pos_neg := db.Query("select id,ad_campaign_id,state_ids,targeting_type_pos from ad_pos_neg")
    if err_ad_pos_neg != nil {
        log.Fatal(err_ad_pos_neg)
    }
    images_sql,err_images := db.Query("select id,image_location,ad_campaign_id from images")
    if err_images != nil {
        log.Fatal(err_images)
    }
    db.Close()

    // for ad_campaign_sql.Next(){
    //     var id string
    //     err = ad_campaign_sql.Scan(&id)
    //     fmt.Println("Positive",id)
    // }
    for state_sql.Next(){
        var id string
        var state_code string
        err = state_sql.Scan(&id,&state_code)
        state_map[id] = state{state_code}
    }
    fmt.Println("States :",state_map)
    for ad_pos_neg_sql.Next(){
        var id string
        var ad_campaign_id int
        var state_ids string
        var state_ids_arr []string
        var targeting_type_pos bool
        err = ad_pos_neg_sql.Scan(&id,&ad_campaign_id,&state_ids,&targeting_type_pos)
        state_ids = state_ids[1:len(state_ids)-1]
        state_ids_arr = strings.Split(state_ids, ",")
        ad_pos_neg_map[id] = ad_pos_neg{ad_campaign_id,state_ids_arr,targeting_type_pos}
    }
    fmt.Println("Ad positive Negative :",ad_pos_neg_map)

    for images_sql.Next(){
        var id string
        var image_location string
        var ad_campaign_id int
        err = images_sql.Scan(&id,&image_location,&ad_campaign_id)
        images_map[id] = images{image_location,ad_campaign_id}
    }
    fmt.Println("Images :",images_map)
    http.HandleFunc("/", getRegion) // set router
    fmt.Println("All Data Loaded")
    fmt.Println("Server Started on port 8069")
    err2 := http.ListenAndServe(":8069", nil) // set listen port
    if err2 != nil {
        log.Fatal("ListenAndServe: ", err2)
    }
}
