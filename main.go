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
  // Vectors to hold the data frame we're about to build
  fmt.Println("names <- c()")
  fmt.Println("AsigBetter <- c()")
  fmt.Println("BsigBetter <- c()")
  for i := 0; i < numSystems; i++ {
    for j := i+1 ; j < numSystems; j++ {
      // First column is the names of the systems under test
      fmt.Print("names <- append(names,paste(\"",path.Base(systems[i]),"\",\"", path.Base(systems[j]),"\",sep=\" & \"))")
      fmt.Println()
      // Next column is "yes" if system A is significantly better than B, "no" otherwise
      fmt.Print("AsigBetter <- append(AsigBetter,if(t.test(",path.Base(systems[i]),"$V1,",path.Base(systems[j]),"$V1,paired=TRUE)$p.value < 0.05) { if(mean(",path.Base(systems[i]),"$V1) > mean(",path.Base(systems[j]),"$V1)) { \"yes\"  } else {\"no\"} } else { \"no\" })")
      fmt.Println()
      // Next column is "yes" if system B is significantly better than A, "no" otherwise
      fmt.Print("BsigBetter <- append(BsigBetter,if(t.test(",path.Base(systems[i]),"$V1,",path.Base(systems[j]),"$V1,paired=TRUE)$p.value < 0.05) { if(mean(",path.Base(systems[i]),"$V1) < mean(",path.Base(systems[j]),"$V1)) { \"yes\"  } else {\"no\"} } else { \"no\" })")
      fmt.Println()
    }
  }
  // Create data frame we've just built
  fmt.Println("table <- data.frame(names,AsigBetter,BsigBetter)")
  // For now, just print the table
  fmt.Print("write.table(table,\"",path.Join(dirName,"mohsen.csv.partial"),"\", row.names=FALSE, col.names=FALSE, sep=\",\")")
  fmt.Println()
}
