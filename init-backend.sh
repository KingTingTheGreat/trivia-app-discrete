cd go-backend &&
(./go-backend.exe || 
    (
        go build && 
        chmod +x go-backend.exe && 
        ./go-backend.exe;
    )
);
