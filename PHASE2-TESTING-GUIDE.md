# Phase 2 Testing Guide: Distributed Session Management

## ✅ Implementation Complete!

Phase 2 has been successfully implemented with distributed session management using Redis. Your backend can now scale horizontally!

**Current Status:**
- ✅ 2 backend replicas running
- ✅ Backend Instance 1: `907a3e63-5271-4712-98a3-a5e64ef1a1a2`
- ✅ Backend Instance 2: `0b0e3fa9-be9c-4789-a988-8938efa62930`
- ✅ Session store using Redis
- ✅ Sticky sessions configured in ingress

---

## 🎯 What's New in Phase 2

### **1. Distributed Session Management**
- Sessions are now stored in Redis (not just in-memory)
- Multiple backend instances can access the same session metadata
- Sessions persist across backend restarts
- Each backend has a unique instance ID

### **2. WebSocket Affinity Tracking**
- Redis tracks which backend instance has the WebSocket for each session
- Prevents WebSocket connection conflicts
- Routes WebSocket requests to the correct backend

### **3. Distributed Locking**
- Redis-based distributed locks prevent race conditions
- Session create/update/delete operations are atomic
- Safe for concurrent access from multiple backends

### **4. Sticky Sessions via Ingress**
- Traefik ingress configured with session affinity
- WebSocket connections stick to the same backend
- Cookie-based session routing

### **5. Backend Instance Heartbeats**
- Each backend sends heartbeats to Redis every 5 seconds
- Failed backends can be detected
- Foundation for future session recovery

---

## 🧪 Testing Instructions

### **Test 1: Verify Multiple Backend Instances**

Check that you have multiple backend pods running:

```bash
kubectl get pods -n gadget-management -l app=gadget-backend
```

**Expected Result:**
```
NAME                      READY   STATUS    RESTARTS   AGE
backend-955747694-m92pj   1/1     Running   0          5m
backend-955747694-sxff5   1/1     Running   0          5m
```

You should see 2 backend pods (or more if you scaled further).

---

### **Test 2: Verify Unique Instance IDs**

Check that each backend has a unique instance ID:

```bash
# Check first backend
kubectl logs -n gadget-management -l app=gadget-backend | grep "Session store initialized"
```

**Expected Result:**
```
2025/11/27 14:30:04 Session store initialized with instance ID: 907a3e63-5271-4712-98a3-a5e64ef1a1a2
2025/11/27 14:30:35 Session store initialized with instance ID: 0b0e3fa9-be9c-4789-a988-8938efa62930
```

Each backend should have a different UUID.

---

### **Test 3: Start a Gadget and Verify Session in Redis**

1. **Start a gadget from the UI:**
   - Open `http://gadget-management.lima.local/`
   - Click "All Gadgets" → "Trace TCP"
   - Select namespace: `demo`
   - Click "Start"

2. **Note the session ID from the UI** (it's the long UUID)

3. **Check Redis for the session:**

```bash
# Connect to Redis
kubectl exec -it -n gadget-management deployment/redis -- redis-cli

# List all session keys
KEYS session:*

# Get a specific session (replace SESSION_ID with your actual session ID)
GET session:<SESSION_ID>

# Check the active sessions index
SMEMBERS sessions:active

# Check which backend has the WebSocket
GET ws:<SESSION_ID>

# Exit Redis
exit
```

**Expected Result:**
- You should see your session key in Redis
- The session data should be a JSON object with session metadata
- The `sessions:active` set should contain your session ID
- The `ws:<SESSION_ID>` key should show which backend instance has the WebSocket

---

### **Test 4: Verify Sessions Appear in Both Backend Instances**

With distributed session management, both backends should see the same sessions:

```bash
# Call the sessions API multiple times
# Due to load balancing, requests will hit different backends
for i in {1..5}; do
  echo "Request $i:"
  curl -s http://gadget-management.lima.local/api/sessions | jq '.[].id'
  echo ""
  sleep 1
done
```

**Expected Result:**
- All requests should return the same session IDs
- Even though requests hit different backend instances, they all see the same sessions from Redis

---

### **Test 5: Test Sticky Sessions for WebSocket**

1. **Open browser DevTools** (Network tab)
2. **Start a gadget** from the UI
3. **Watch the WebSocket connection** in Network tab:
   - Look for a connection to `/ws/<session-id>`
   - Check the "Headers" tab
   - Look for the `Cookie` header with `backend-session=...`

4. **Refresh the page** and start another gadget
5. **Verify** that the WebSocket stays connected to the same backend

**Expected Result:**
- The `backend-session` cookie should be set
- WebSocket connections should stick to the same backend instance
- No "Session is on a different backend instance" errors

---

### **Test 6: Test Session Persistence Across Backend Restart**

1. **Start a gadget** and note the session ID

2. **Check which backend has the session:**
```bash
kubectl logs -n gadget-management -l app=gadget-backend --tail=50 | grep "Started gadget"
```

3. **Delete one backend pod** (Kubernetes will recreate it):
```bash
kubectl delete pod -n gadget-management <backend-pod-name>
```

4. **Wait for the new pod to start:**
```bash
kubectl get pods -n gadget-management -l app=gadget-backend -w
```

5. **Check if the session is still listed:**
```bash
curl http://gadget-management.lima.local/api/sessions | jq
```

**Expected Result:**
- Sessions started on the deleted backend will show in the list (metadata is in Redis)
- However, the actual gadget process is gone, so it won't produce new events
- This demonstrates session metadata persistence (not full session recovery yet - that's for a future phase)

---

### **Test 7: Scale Backend to 3 Replicas**

```bash
kubectl scale deployment backend -n gadget-management --replicas=3
kubectl get pods -n gadget-management -l app=gadget-backend -w
```

Wait for all 3 pods to be running, then:

```bash
# Check that all 3 backends registered themselves
kubectl exec -it -n gadget-management deployment/redis -- redis-cli KEYS backend:*:heartbeat
```

**Expected Result:**
- You should see 3 backend pods running
- Redis should have 3 heartbeat keys (one per backend instance)

---

### **Test 8: Verify Backend Heartbeats**

```bash
# Connect to Redis
kubectl exec -it -n gadget-management deployment/redis -- redis-cli

# List all backend heartbeats
KEYS backend:*:heartbeat

# Check a specific backend's last heartbeat timestamp
GET backend:<instance-id>:heartbeat

# Exit
exit
```

**Expected Result:**
- Each backend should have a heartbeat key
- The value should be a recent Unix timestamp
- Heartbeats should update every 5 seconds

---

### **Test 9: Distributed Locking Test**

This test verifies that distributed locks prevent race conditions:

```bash
# Start multiple gadgets simultaneously from the UI
# (Open multiple browser tabs and click "Start" at the same time)

# Then check Redis for lock keys
kubectl exec -it -n gadget-management deployment/redis -- redis-cli KEYS "lock:session:*"
```

**Expected Result:**
- You should see temporary lock keys during session creation
- Locks should be released after the operation completes
- No session corruption or conflicts should occur

---

### **Test 10: Historical Data Still Works**

With Phase 2, historical data persistence (Phase 1) should still work:

1. **Start a gadget** and let it run for 1-2 minutes
2. **Stop the gadget**
3. **Wait 30 seconds** for events to be processed
4. **Go to History tab** and search
5. **Try session replay**

**Expected Result:**
- ✅ All Phase 1 features still work
- ✅ Events are persisted to TimescaleDB
- ✅ History search works
- ✅ Session replay works

---

## 🔍 Monitoring and Debugging

### **Check Redis Session Data**

```bash
# Connect to Redis
kubectl exec -it -n gadget-management deployment/redis -- redis-cli

# Get all session-related keys
KEYS *

# Count active sessions
SCARD sessions:active

# Get backend sessions for a specific instance
SMEMBERS backend:<instance-id>:sessions

# Get all WebSocket registrations
KEYS ws:*
```

### **Check Backend Logs**

```bash
# All backends
kubectl logs -n gadget-management -l app=gadget-backend --tail=100

# Specific backend
kubectl logs -n gadget-management <backend-pod-name> --tail=50

# Follow logs
kubectl logs -n gadget-management -l app=gadget-backend -f
```

### **Check Ingress Routing**

```bash
# Check ingress configuration
kubectl describe ingress gadget-management -n gadget-management

# Look for sticky session annotations
kubectl get ingress gadget-management -n gadget-management -o yaml | grep affinity
```

### **Monitor Backend Replicas**

```bash
# Watch backend scaling
kubectl get deployment backend -n gadget-management -w

# Check backend service endpoints
kubectl get endpoints backend -n gadget-management

# Check which pods are behind the service
kubectl get pods -n gadget-management -l app=gadget-backend -o wide
```

---

## 📊 Architecture Overview

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Browser    │     │  Browser    │     │  Browser    │
│  Client 1   │     │  Client 2   │     │  Client 3   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Ingress   │
                    │  (Traefik)  │
                    │  Sticky     │
                    │  Sessions   │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
     ┌─────▼─────┐   ┌─────▼─────┐  ┌─────▼─────┐
     │Backend-1  │   │Backend-2  │  │Backend-3  │
     │ ID: 907.. │   │ ID: 0b0.. │  │ ID: xxx.. │
     └─────┬─────┘   └─────┬─────┘  └─────┬─────┘
           │               │               │
           └───────────────┼───────────────┘
                           │
                    ┌──────▼──────┐
                    │    Redis    │
                    │             │
                    │ • Sessions  │
                    │ • WS Track  │
                    │ • Heartbeat │
                    │ • Locks     │
                    └─────────────┘
```

**Key Components:**
- **Ingress**: Routes requests with sticky sessions (cookie-based)
- **Backend Replicas**: Multiple instances, each with unique ID
- **Redis**: Central store for session metadata, WebSocket tracking, locks, and heartbeats
- **Gadget Sessions**: Local to each backend (where the kubectl-gadget process runs)

---

## 🎉 Success Criteria

**Phase 2 is FULLY COMPLETE when:**
- ✅ Multiple backend replicas running
- ✅ Each backend has unique instance ID
- ✅ Sessions stored in Redis
- ✅ WebSocket connections tracked
- ✅ Sticky sessions working (no 502 errors)
- ✅ Sessions persist across backend restarts (metadata only)
- ✅ All Phase 1 features still work
- ✅ Backend heartbeats visible in Redis
- ✅ Distributed locks prevent race conditions

---

## 🚀 Scaling Operations

### **Scale Up**

```bash
# Scale to 5 replicas
kubectl scale deployment backend -n gadget-management --replicas=5

# Watch the rollout
kubectl rollout status deployment backend -n gadget-management

# Verify all instances registered
kubectl exec -it -n gadget-management deployment/redis -- redis-cli KEYS backend:*:heartbeat
```

### **Scale Down**

```bash
# Scale back to 2 replicas
kubectl scale deployment backend -n gadget-management --replicas=2

# Watch pods terminate
kubectl get pods -n gadget-management -l app=gadget-backend -w
```

### **Auto-Scaling (Future)**

You can enable Horizontal Pod Autoscaler:

```bash
kubectl autoscale deployment backend -n gadget-management \
  --cpu-percent=70 \
  --min=2 \
  --max=10
```

---

## 🐛 Troubleshooting

### **Issue: "Session is on a different backend instance" Error**

**Cause:** WebSocket trying to connect to a backend that doesn't have the gadget session.

**Solution:**
- Check if sticky sessions are working: `kubectl describe ingress gadget-management -n gadget-management`
- Verify the cookie is being set in browser DevTools
- Check if the ingress has the affinity annotations

### **Issue: Sessions Disappear After Backend Restart**

**Cause:** Session metadata is in Redis, but the actual gadget process is on a specific backend.

**Expected Behavior:** This is normal for Phase 2. Session *metadata* persists, but the actual running gadget process doesn't (full session recovery is a future enhancement).

### **Issue: Redis Connection Errors**

```bash
# Check if Redis is running
kubectl get pods -n gadget-management -l app=redis

# Check Redis logs
kubectl logs -n gadget-management deployment/redis --tail=50

# Test Redis connection
kubectl exec -it -n gadget-management deployment/redis -- redis-cli PING
```

---

## 📈 What's Next?

With Phase 2 complete, you now have:
✅ Horizontal scaling capability
✅ Session persistence across restarts
✅ High availability
✅ Load balancing

**Possible Future Enhancements:**
- **Session Recovery**: Reconnect to running gadgets after backend restart
- **Auto-Scaling**: HPA based on CPU/memory or custom metrics
- **Session Migration**: Move sessions between backends
- **Advanced Load Balancing**: Route based on session workload
- **Multi-Region**: Deploy across multiple regions/clusters

Congratulations on completing Phase 2! 🎉
