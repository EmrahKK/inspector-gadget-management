# Testing Guide: Historical Features & Session Replay

## ✅ Implementation Complete!

The frontend has been successfully updated with historical data viewing and session replay features.

**New Build:** `index-m4A53HQ9.js`

---

## 🎯 New Features

### 1. **History Tab** - Query and analyze historical events
- Date/time range filtering
- Event type filtering (trace_sni, trace_tcp, snapshot_process, snapshot_socket)
- Namespace filtering
- Session ID search
- Quick time range buttons (Last Hour, Last 24h, Last 7 days)
- Export to JSON
- Event expansion to see full data
- Direct replay from any event

### 2. **Session Replay** - Replay historical sessions
- Load all events for a session
- Play/pause controls
- Adjustable playback speed (0.5x, 1x, 2x, 5x, 10x)
- Timeline scrubbing
- Event-by-event navigation
- Session statistics
- Export session events to JSON

---

## 🧪 Testing Instructions

### **Prerequisites: Generate Some Historical Data**

Before testing the historical features, you need to generate some events. Here's how:

1. **Open the UI** in your browser:
   ```
   http://gadget-management.lima.local/
   ```

2. **Start a gadget** (e.g., trace_tcp in the demo namespace)
   - Click "All Gadgets" → "Trace TCP"
   - Select namespace: `demo`
   - Configure any filters you want
   - Click "Start"

3. **Generate some traffic** in your cluster to create events
   - If using trace_tcp, any network activity in the demo namespace will create events
   - For snapshot gadgets, they'll capture data immediately

4. **Let it run for 1-2 minutes** to collect events

5. **Stop the gadget** - This is important! The events need to be persisted to the database

6. **Wait 5-10 seconds** for the consumer worker to process and write events to TimescaleDB

---

## 📋 Test Cases

### **Test 1: View History Tab**

1. Click **"History"** in the left sidebar (at the bottom of the catalog section)
2. You should see the History view with filters

**Expected Result:**
- ✅ Clean UI with filter panel
- ✅ Event type dropdown with all gadget types
- ✅ Namespace input field
- ✅ Session ID input field
- ✅ Start/End time pickers
- ✅ Limit selector (50/100/500/1000)
- ✅ Quick range buttons

---

### **Test 2: Query Last Hour of Events**

1. In History view, click **"Last Hour"** button
2. Click **"Search"** button

**Expected Result:**
- ✅ Events appear in a list
- ✅ Each event shows:
  - Event type badge
  - Timestamp
  - Preview data (process name, pod, namespace, etc.)
  - "Replay Session" link
  - Expand/collapse button for full data
- ✅ Header shows: "Results: X events from Y sessions"
- ✅ No errors

**If No Events:**
- Make sure you completed the prerequisites (started a gadget and let it run)
- Try a longer time range (Last 24h)
- Check backend logs: `kubectl logs -n gadget-management -l app=gadget-backend --tail=50`

---

### **Test 3: Filter by Event Type**

1. Select event type from dropdown (e.g., "Trace TCP")
2. Click "Search"

**Expected Result:**
- ✅ Only events matching that type appear
- ✅ Results update correctly

---

### **Test 4: Filter by Namespace**

1. Enter namespace in the field (e.g., "demo")
2. Click "Search"

**Expected Result:**
- ✅ Only events from that namespace appear

---

### **Test 5: Search by Session ID**

1. Copy a session ID from the events list (the long UUID)
2. Paste it into the "Session ID" field
3. Click "Search"

**Expected Result:**
- ✅ Only events from that specific session appear
- ✅ Results should be filtered to one session

---

### **Test 6: Export Events to JSON**

1. After searching, click **"Export JSON"** button (top right of results)

**Expected Result:**
- ✅ JSON file downloads automatically
- ✅ Filename: `gadget-events-[timestamp].json`
- ✅ File contains array of event objects

---

### **Test 7: Expand Event Details**

1. Find an event in the results
2. Click **"Show full data"**

**Expected Result:**
- ✅ Raw JSON expands below the event
- ✅ Data is formatted and readable
- ✅ Click "Hide full data" collapses it

---

### **Test 8: Launch Session Replay from History**

1. Find an event in the results
2. Click the **"Replay Session"** link (blue underlined text next to timestamp)

**Expected Result:**
- ✅ Session Replay modal opens
- ✅ Modal shows session ID at top
- ✅ Statistics panel shows: Type, Namespace, Duration, Total Events
- ✅ Playback controls visible (Reset, Play/Pause, scrubber, speed selector, Export)
- ✅ First event displayed in "Current Event" section
- ✅ Timeline shows all events on the right side

---

### **Test 9: Session Replay - Playback Controls**

**With Session Replay modal open:**

1. **Test Play/Pause:**
   - Click ▶ Play button
   - Events should auto-advance
   - Click ⏸ Pause to stop

2. **Test Speed Control:**
   - Change speed to 2x
   - Click Play
   - Events should advance faster

3. **Test Reset:**
   - Click ↻ Reset button
   - Should jump back to first event

4. **Test Scrubbing:**
   - Drag the progress bar slider
   - Should jump to different events

5. **Test Timeline Navigation:**
   - Click any event in the timeline (right panel)
   - Should jump to that event

**Expected Results:**
- ✅ All controls work smoothly
- ✅ Current event updates correctly
- ✅ Progress bar matches current position
- ✅ No lag or freezing

---

### **Test 10: Session Replay - Export**

1. In Session Replay modal, click **"Export"** button

**Expected Result:**
- ✅ JSON file downloads
- ✅ Filename: `session-[session-id]-events.json`
- ✅ Contains all events for that session

---

### **Test 11: Session Replay - Event Details**

1. In Session Replay modal, look at "Current Event" section
2. Event data should be displayed as key-value pairs
3. Click **"Show raw JSON"** to expand

**Expected Result:**
- ✅ Key fields displayed nicely (comm, pod, namespace, etc.)
- ✅ Raw JSON available in expandable section
- ✅ All data readable

---

### **Test 12: Close and Reopen**

1. Close Session Replay modal (X button)
2. Go back to History
3. Search again and click another "Replay Session"

**Expected Result:**
- ✅ Modal closes cleanly
- ✅ Can reopen with different session
- ✅ Previous session data doesn't persist

---

### **Test 13: Combined Filters**

1. In History view:
   - Select event type: "Trace TCP"
   - Enter namespace: "demo"
   - Set time range: Last 24h
   - Set limit: 50
2. Click "Search"

**Expected Result:**
- ✅ Results match ALL filters
- ✅ Max 50 events returned
- ✅ All are trace_tcp type from demo namespace

---

### **Test 14: Reset Filters**

1. Set some filters
2. Click **"Reset"** button

**Expected Result:**
- ✅ All filters clear
- ✅ Results list clears
- ✅ Ready for new search

---

### **Test 15: Empty State**

1. Set filters that won't match anything (e.g., namespace "nonexistent")
2. Click "Search"

**Expected Result:**
- ✅ No error message
- ✅ Shows empty state: "No events found"
- ✅ Helpful message: "Adjust your filters and click Search..."

---

## 🐛 Troubleshooting

### **No Events Appear in History**

**Check 1: Are events being persisted?**
```bash
# Connect to TimescaleDB
kubectl exec -it -n gadget-management deployment/timescaledb -- \
  psql -U gadget -d gadget_events -c "SELECT COUNT(*) FROM gadget_events;"

# Should return a count > 0
```

**Check 2: Are events in Redis?**
```bash
# Check Redis stream
kubectl exec -n gadget-management deployment/redis -- redis-cli XLEN gadget:events

# If 0, events haven't been published
# If >0, events are being published but not consumed
```

**Check 3: Backend logs**
```bash
kubectl logs -n gadget-management -l app=gadget-backend --tail=50

# Should see:
# - "Connected to Redis"
# - "Connected to PostgreSQL"
# - "Starting event consumer..."
# - No errors about publishing or consuming
```

**Check 4: Was WebSocket connected?**
Events are only published when a WebSocket is active. Make sure you:
1. Started a gadget from the UI (not just curl)
2. Kept the gadget view open for a bit
3. Saw events streaming in real-time

---

### **Session Replay Shows "No events found"**

This means the session exists but has no events in the database. Try:
1. Start a new gadget and let it run for 1-2 minutes
2. Stop it
3. Wait 10 seconds for events to be written
4. Try replaying that session

---

### **Playback is Jerky or Laggy**

This is normal if there are many events. Try:
1. Reduce playback speed to 0.5x
2. Use the timeline to jump directly to events
3. Export and analyze the JSON instead

---

### **API Errors**

If you see errors like "Failed to query events":
```bash
# Check API is working
curl "http://gadget-management.lima.local/api/events?limit=10"

# Should return JSON array (possibly empty)
# If 404 or 500, check backend logs
```

---

## 📊 Expected Data Flow

```
1. Start Gadget in UI
   ↓
2. WebSocket connects to backend
   ↓
3. Gadget generates events
   ↓
4. Backend publishes to Redis Streams
   ↓
5. Backend forwards to WebSocket (real-time view)
   ↓
6. Consumer worker reads from Redis
   ↓
7. Consumer writes to TimescaleDB
   ↓
8. Events available in History view!
```

---

## ✨ Tips for Best Testing Experience

1. **Generate varied data:**
   - Run multiple gadget types (trace_tcp, trace_sni, snapshot_process)
   - Use different namespaces
   - Run for different durations

2. **Create multiple sessions:**
   - Start and stop several gadgets
   - This lets you test filtering and replay better

3. **Test with real traffic:**
   - trace_tcp works best with active network connections
   - trace_sni needs TLS connections
   - snapshot gadgets work immediately

4. **Check database growth:**
   ```bash
   # Monitor event count
   watch -n 5 'kubectl exec -n gadget-management deployment/timescaledb -- \
     psql -U gadget -d gadget_events -c "SELECT COUNT(*) FROM gadget_events;"'
   ```

---

## 🎉 Success Criteria

**Phase 1 is FULLY COMPLETE when:**
- ✅ History tab appears in sidebar
- ✅ Can query events with filters
- ✅ Can see event details
- ✅ Can export events to JSON
- ✅ Can click "Replay Session"
- ✅ Session Replay modal opens
- ✅ Can play/pause/scrub through events
- ✅ Can export session data
- ✅ All controls work smoothly
- ✅ No console errors

---

## 🚀 What's Next?

After testing, you now have:
1. ✅ Full data persistence
2. ✅ Historical querying
3. ✅ Session replay
4. ✅ Export capabilities

**Ready for Phase 2?**
- Distributed session management
- Horizontal backend scaling
- Redis-based session state
- Load balancing

Enjoy exploring your gadget history! 🎯
