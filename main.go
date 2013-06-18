package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "path"
  "log"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println(os.Args[0], ": generates an R script for producing a file in Mohsen's format")
    fmt.Println("Usage:")
    fmt.Println("   ",os.Args[0], " <path to directory containing .eval files>")
    return
  }

  // Read the collection CSVs
  dirName := os.Args[1]
  fileListing, err:= ioutil.ReadDir(dirName);
  if err != nil {
    log.Fatal(err)
  }
  systems := make([]string,0,len(fileListing))
  for _,file := range fileListing {
     if path.Ext(file.Name()) == ".eval"{
        systems = append(systems,path.Join(dirName,file.Name()))
     }
  }



  numSystems := len(systems)
  // First, we need R to read in the data
  for i := 0; i < numSystems; i++ {
    fmt.Print(path.Base(systems[i]), " <- read.csv(sep=\" \",header=FALSE,file=\"",systems[i],"\")")
    fmt.Println()
  }
  fmt.Println("names <- c()")
  fmt.Println("AsigBetter <- c()")
  fmt.Println("BsigBetter <- c()")
  for i := 0; i < numSystems; i++ {
    for j := i+1 ; j < numSystems; j++ {
      fmt.Print("names <- append(names,paste(\"",path.Base(systems[i]),"\",\"", path.Base(systems[j]),"\",sep=\" & \"))")
      fmt.Println()
      fmt.Print("AsigBetter <- append(AsigBetter,if(t.test(",path.Base(systems[i]),"$V2,",path.Base(systems[j]),"$V2,paired=TRUE)$p.value < 0.05) { if(mean(",path.Base(systems[i]),"$V2) > mean(",path.Base(systems[j]),"$V2)) { \"yes\"  } else {\"no\"} } else { \"no\" })")
      fmt.Println()
      fmt.Print("BsigBetter <- append(BsigBetter,if(t.test(",path.Base(systems[i]),"$V2,",path.Base(systems[j]),"$V2,paired=TRUE)$p.value < 0.05) { if(mean(",path.Base(systems[i]),"$V2) < mean(",path.Base(systems[j]),"$V2)) { \"yes\"  } else {\"no\"} } else { \"no\" })")
      fmt.Println()
    }
  }
  fmt.Println("table <- data.frame(names,AsigBetter,BsigBetter)")
  fmt.Println("table");
}
