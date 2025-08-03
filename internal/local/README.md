# intenal/local
localに関わるものは、db,file,webとあるが、それらに属しないlocalを管理する。  
簡単に思いつくものとしては、local osの機能を使うlocalや、object storageを扱う機能。  
object storageはそれぞれのfile形式があり、それらはfile以下で取り扱うが、純粋にobject storageにアクセスする機能はこちらで管理したい。  
