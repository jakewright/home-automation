# Run

Run is a tool for running home-automation services locally. 

## Package structure

                                            
```                                            
       +-->service <--+                     
       |              |                     
       |              |                     
       v              v                     
    golang         compose     } run systems
       ^              ^                     
       |              |                     
       v              |                     
    docker            |                     
       ^              |                     
       |              |                     
       v              v                     
  +------------------------+                
  |         Docker         |                
  +------------------------+                
```

https://textik.com/#827c8a3c02be9468
