run: build
	cd nextjs-frontend && npm start & cd go-backend && ./go-backend.exe

build:
	if [ ! -d nextjs-frontend/.next ]; then cd nextjs-frontend && npm run build; fi;
	if [ ! -f go-backend/go-backend.exe ]; then cd go-backend && go build; fi;

clean:
	rm -r nextjs-frontend/.next 
	rm go-backend/go-backend.exe
