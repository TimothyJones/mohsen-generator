package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "path"
  "log"
)

func main() {
  noZeros := 0
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
    fmt.Print(path.Base(systems[i]), " <- ", path.Base(systems[i]),"$V1[1:50]")
    fmt.Println()
  }
  //fmt.Println("exactRankTests")
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
      // Let's try for no zeroes
      fmt.Print("frame <- data.frame(",path.Base(systems[i]),",",path.Base(systems[j]),")")
      fmt.Println()
      if noZeros == 1 {
        fmt.Print("frame <- subset(frame, ",path.Base(systems[i])," != 0 & ",path.Base(systems[j])," != 0)")
        fmt.Println()
      }
      fmt.Print("AsigBetter <- append(AsigBetter,if(t.test(frame$",path.Base(systems[i]),",frame$",path.Base(systems[j]),",paired=TRUE)$p.value < 0.05) { if(mean(frame$",path.Base(systems[i]),") > mean(frame$",path.Base(systems[j]),")) { \"yes\"  } else {\"no\"} } else { \"no\" })")
//      fmt.Print("AsigBetter <- append(AsigBetter,if(perm.test(",path.Base(systems[i]),",",path.Base(systems[j]),",paired=TRUE)$p < 0.05) { if(mean(",path.Base(systems[i]),") > mean(",path.Base(systems[j]),")) { \"yes\"  } else {\"no\"} } else { \"no\" })")


      fmt.Println()
      // Next column is "yes" if system B is significantly better than A, "no" otherwise
     fmt.Print("BsigBetter <- append(BsigBetter,if(t.test(frame$",path.Base(systems[i]),",frame$",path.Base(systems[j]),",paired=TRUE)$p.value < 0.05) { if(mean(frame$",path.Base(systems[i]),") < mean(frame$",path.Base(systems[j]),")) { \"yes\"  } else {\"no\"} } else { \"no\" })")
 //     fmt.Print("BsigBetter <- append(BsigBetter,if(perm.test(",path.Base(systems[i]),",",path.Base(systems[j]),",paired=TRUE)$p < 0.05) { if(mean(",path.Base(systems[i]),") < mean(",path.Base(systems[j]),")) { \"yes\"  } else {\"no\"} } else { \"no\" })")
      fmt.Println()
    }
  }
  // Create data frame we've just built
  fmt.Println("table <- data.frame(names,AsigBetter,BsigBetter)")
  // For now, just print the table
  fmt.Print("write.table(table,\"",path.Join(dirName,"mohsen.csv.partial"),"\", row.names=FALSE, col.names=FALSE, sep=\",\")")
  fmt.Println()
}
