# Troubleshooting Guide

## Common Issues and Solutions

### 1. Backend Issues

#### Port 8080 Already in Use
**Error:** `bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in backend/cmd/server/main.go
http.ListenAndServe(":8081", handler)
```

#### Go Module Errors
**Error:** `package not found` or `module not found`

**Solution:**
```bash
cd backend
go mod tidy
go mod download
```

#### CORS Errors
**Error:** `CORS policy: No 'Access-Control-Allow-Origin' header`

**Solution:**
Check `backend/cmd/server/main.go` CORS configuration:
```go
AllowedOrigins: []string{"http://localhost:5173"}
```

### 2. Frontend Issues

#### Port 5173 Already in Use
**Error:** `Port 5173 is in use`

**Solution:**
Edit `frontend/vite.config.js`:
```js
server: {
  port: 5174
}
```

#### Module Not Found
**Error:** `Cannot find module` or `Module not found`

**Solution:**
```bash
cd frontend
rm -rf node_modules
rm package-lock.json
npm install
```

#### API Connection Failed
**Error:** `Failed to fetch` or `Network error`

**Solution:**
1. Check backend is running on port 8080
2. Check browser console for errors
3. Verify API URL in `frontend/src/lib/api/client.js`

### 3. Build Issues

#### Backend Build Fails
**Error:** Compilation errors

**Solution:**
```bash
cd backend
go clean
go build -v ./cmd/server/main.go
```

#### Frontend Build Fails
**Error:** Build errors

**Solution:**
```bash
cd frontend
npm run build
```

### 4. Runtime Issues

#### Empty Data on Frontend
**Problem:** Pages show loading but no data

**Solution:**
1. Check backend is running
2. Open browser DevTools > Network tab
3. Check API responses
4. Verify CORS is working

#### Cart Not Working
**Problem:** Items not adding to cart

**Solution:**
1. Check browser console for errors
2. Verify cart store in `frontend/src/lib/stores/cart.js`
3. Clear browser cache and reload

#### Orders Not Creating
**Problem:** Checkout fails

**Solution:**
1. Check backend logs
2. Verify request payload in Network tab
3. Check order structure matches backend model

### 5. Development Issues

#### Hot Reload Not Working

**Backend:**
Use `air` for hot reload:
```bash
go install github.com/cosmtrek/air@latest
cd backend
air
```

**Frontend:**
Restart dev server:
```bash
cd frontend
npm run dev
```

#### Changes Not Reflecting

**Solution:**
1. Hard refresh browser (Cmd+Shift+R on Mac, Ctrl+Shift+R on Windows)
2. Clear browser cache
3. Restart dev servers

### 6. Dependency Issues

#### Go Dependencies
```bash
cd backend
go clean -modcache
go mod download
```

#### Node Dependencies
```bash
cd frontend
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### 7. Permission Issues

#### Script Not Executable
**Error:** `Permission denied: ./start.sh`

**Solution:**
```bash
chmod +x start.sh
```

### 8. Database Issues (Future)

When you add database:

#### Connection Failed
- Check database is running
- Verify connection string
- Check credentials

#### Migration Errors
- Run migrations manually
- Check migration files
- Verify database schema

## Debugging Tips

### Backend Debugging

1. **Add Logging:**
```go
log.Printf("Debug: %+v", variable)
```

2. **Check Request:**
```go
body, _ := ioutil.ReadAll(r.Body)
log.Printf("Request body: %s", body)
```

3. **Use Postman/Insomnia:**
Test API endpoints directly

### Frontend Debugging

1. **Console Logging:**
```javascript
console.log('Debug:', data);
```

2. **Svelte DevTools:**
Install Svelte DevTools browser extension

3. **Network Tab:**
Check all API requests and responses

## Performance Issues

### Backend Slow
- Check handler logic
- Add caching if needed
- Profile with `pprof`

### Frontend Slow
- Check bundle size: `npm run build`
- Optimize images
- Lazy load components

## Common Mistakes

1. **Forgetting to start backend** before frontend
2. **Wrong API URL** in frontend config
3. **CORS not configured** properly
4. **Port conflicts** with other services
5. **Missing dependencies** after git clone

## Getting Help

If issues persist:

1. Check browser console for errors
2. Check backend terminal for logs
3. Verify all dependencies are installed
4. Try restarting both servers
5. Clear all caches and rebuild

## Useful Commands

```bash
# Check Go version
go version

# Check Node version
node --version

# Check if port is in use
lsof -i :8080
lsof -i :5173

# Kill process on port
kill -9 $(lsof -t -i:8080)

# Clean everything and restart
cd backend && go clean && cd ../frontend && rm -rf node_modules && npm install
```

## Environment Variables (Future)

When adding .env files:

```bash
# Backend .env
PORT=8080
DATABASE_URL=...

# Frontend .env
VITE_API_URL=http://localhost:8080
```

## Testing Checklist

Before reporting issues, verify:

- [ ] Backend is running on port 8080
- [ ] Frontend is running on port 5173
- [ ] No console errors in browser
- [ ] No errors in backend terminal
- [ ] Dependencies are installed
- [ ] Ports are not in use by other apps
- [ ] CORS is configured correctly
- [ ] API endpoints are accessible

## Quick Reset

If everything is broken:

```bash
# Stop all servers (Ctrl+C)

# Backend
cd backend
go clean
rm -rf go.sum
go mod tidy
go mod download

# Frontend
cd ../frontend
rm -rf node_modules package-lock.json .svelte-kit
npm install

# Restart
cd ..
./start.sh
```

---

Still having issues? Check the documentation files or create an issue with:
- Error message
- Steps to reproduce
- Environment details (OS, Go version, Node version)
