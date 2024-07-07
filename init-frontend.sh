cd nextjs-frontend && 
(npm start ||
    (
        npm install && 
        npm run build && 
        npm start
    )
);