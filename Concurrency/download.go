package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"                       
	"os"
	"time"
)
              
func main() {

	images:=[10]string{
		"https://static.wikia.nocookie.net/deathnote/images/7/76/DEATH_NOTE_anime.jpg/revision/latest/top-crop/width/360/height/450?cb=20170720215429",
		"https://i0.wp.com/livewire.thewire.in/wp-content/uploads/2020/08/Screen-Shot-2020-08-27-at-2.20.06-PM.png?resize=743%2C384&ssl=1",		
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSkX4QB2ZdvRNTCfGwTWaeEKvfXIBmD2pPEJA&usqp=CAU",
		"https://static.wikia.nocookie.net/p__/images/0/0c/L_Lawliet.png/revision/latest?cb=20210120194215&path-prefix=protagonist",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT6_Rhair-aCW6_bAwiqRE6sYvsuYNOebriyA&usqp=CAU",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRO3cqs1nIZ7mZffD8i8uYen20ZEIvHxLhfjQ&usqp=CAU",
		"https://pbs.twimg.com/profile_images/719653899594964994/MNyI_8SQ_400x400.jpg",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQhoUlzKM6EsLydXWtWGdY-Lfr5j0SHmBiEnw&usqp=CAU",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ-a8ed0mRWvQXCARVLF_uOXzoW9HS0ANzXLw&usqp=CAU",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTluSqr7D0n-VBsGraZ5Ai_V6mqGDhcvzvDeg&usqp=CAU"}

	channelCapacity:=10

	ch:=make(chan bool,channelCapacity)

	for i:=0;i<10;i++{
		go image_download(i,ch,images)
	}	
	for i:=0;i<10;i++{
		<-ch
	}	
	fmt.Printf("File downlaod in current working directory")
}

func image_download(i int,ch chan bool,images [10]string){

	time.Sleep(time.Second*2)
	var fileName string
	if i<5{
		fileName = fmt.Sprint("L",i,".png")
	}else{
		fileName = fmt.Sprint("Captain_Levi",i,".png")
	}
	URL := images[i]
	err := downloadFile(URL,"images/"+fileName)
	if err != nil {
		log.Fatal(err)
	}
	ch<-true
	fmt.Println(i,"downloaded")

}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}      