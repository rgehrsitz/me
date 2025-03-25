@echo off
echo Building AI-Powered Personal Knowledge Base...

echo Building frontend...
cd web
echo Cleaning node_modules to avoid file lock issues...
if exist node_modules (
    rmdir /s /q node_modules
    echo Old node_modules removed
)

echo Installing global dependencies...
call npm install -g vite
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Global vite installation failed, will try to use local install...
)

echo Creating .npmrc to disable husky...
echo ignore-scripts=true > .npmrc
echo CI=true >> .npmrc

echo Installing npm dependencies...
call npm install --no-audit --no-fund
if %ERRORLEVEL% NEQ 0 (
    echo Warning: npm install encountered issues but continuing build process...
)

echo Removing .npmrc...
del .npmrc

echo Checking for index.html entry point...
if not exist public\index.html (
    echo Error: Could not find public/index.html entry point!
    exit /b 1
)

echo Creating vite.config.js...
echo import { defineConfig } from 'vite'; > vite.config.build.js
echo import { svelte } from '@sveltejs/vite-plugin-svelte'; >> vite.config.build.js
echo. >> vite.config.build.js
echo export default defineConfig({ >> vite.config.build.js
echo   plugins: [svelte()], >> vite.config.build.js
echo   root: './', >> vite.config.build.js
echo   publicDir: 'public', >> vite.config.build.js
echo   build: { >> vite.config.build.js
echo     outDir: 'dist', >> vite.config.build.js
echo     emptyOutDir: true, >> vite.config.build.js
echo     assetsDir: 'assets', >> vite.config.build.js
echo   } >> vite.config.build.js
echo }); >> vite.config.build.js

echo Building frontend application...
call npx vite build --config vite.config.build.js
if %ERRORLEVEL% NEQ 0 (
    echo Error: Frontend build failed! Trying with full config...
    call npx vite build -f --root . --outDir dist
    if %ERRORLEVEL% NEQ 0 (
        echo Error: All build attempts failed!
        exit /b %ERRORLEVEL%
    )
)

echo Cleaning up temporary files...
del vite.config.build.js

cd ..

echo Building backend...
echo Checking for C compiler...
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Warning: GCC compiler not found. Attempting to build without CGO...
    echo To fix this permanently, install MinGW-w64 and add it to your PATH.
    set CGO_ENABLED=0
) else (
    echo C compiler found, building with CGO enabled...
    set CGO_ENABLED=1
)

if not exist bin (
    mkdir bin
)

go build -o bin/pkb.exe cmd/server/main.go
if %ERRORLEVEL% NEQ 0 (
    echo Trying alternative build method...
    set CGO_ENABLED=0
    go build -o bin/pkb.exe cmd/server/main.go
    if %ERRORLEVEL% NEQ 0 (
        echo Error: Backend build failed!
        exit /b %ERRORLEVEL%
    ) else (
        echo Built successfully with CGO disabled.
    )
) else (
    echo Backend built successfully!
)

echo Done! Run bin/pkb.exe to start the server.
