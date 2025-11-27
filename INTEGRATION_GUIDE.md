# Integration Guide - New UI Implementation

## What Was Done

I've successfully integrated the new catalog-based UI design from `new-ui-ideas/` into your existing codebase. Here's a comprehensive overview:

### 1. Dependencies Added

Updated `frontend/package.json` with:
- **lucide-react**: ^0.294.0 - Icon library for the UI
- **recharts**: ^2.10.3 - Data visualization library for charts
- **tailwindcss**: ^3.4.0 - CSS framework for styling
- **postcss**: ^8.4.32 - CSS processor
- **autoprefixer**: ^10.4.16 - CSS autoprefixer

### 2. Configuration Files Created

- `frontend/tailwind.config.js` - Tailwind CSS configuration
- `frontend/postcss.config.js` - PostCSS configuration
- `frontend/src/index.css` - Updated with Tailwind directives

### 3. New Components Created

#### `frontend/src/components/GadgetCard.tsx`
A reusable card component for displaying gadgets in the catalog with:
- Icon and category badge
- Running state indicator
- Hover effects
- Click handler for opening the runner

#### `frontend/src/components/Runner.tsx`
The main gadget runner interface with:
- Configuration panel (namespace, pod name, TCP filters)
- Start/Stop controls
- Active configuration display
- Visual and Raw JSON output tabs
- Export functionality
- Integration with existing backend APIs

### 4. Updated App Component

`frontend/src/App.tsx` now features:
- **Catalog View**: Card-based gadget browser organized by categories
- **Sidebar Navigation**:
  - Active Operations panel showing running gadgets
  - Category filters (All, Trace, Top, Snapshot, Profile, Audit)
- **Runner View**: Full-screen runner interface when a gadget is selected
- **Global State Management**: Tracks running sessions across the app
- **Smart Resume**: When clicking a running gadget, it shows existing output
- **Visual Indicators**: Running badges, pulse animations

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Browser UI                            │
│  ┌──────────────┐  ┌──────────────────────────────┐    │
│  │   Sidebar    │  │      Main Content            │    │
│  │              │  │                              │    │
│  │ • Active Ops │  │ Catalog View:                │    │
│  │ • Categories │  │  - Gadget Cards (3 gadgets)  │    │
│  │              │  │                              │    │
│  │              │  │ Runner View:                 │    │
│  │              │  │  - Configuration Panel       │    │
│  │              │  │  - Output Display            │    │
│  │              │  │  - Visual/Raw JSON Tabs      │    │
│  └──────────────┘  └──────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                          ↓
                REST API + WebSocket
                          ↓
┌─────────────────────────────────────────────────────────┐
│               Go Backend (No Changes)                    │
│  • Session Management                                    │
│  • WebSocket Streaming                                   │
│  • kubectl-gadget Execution                             │
└─────────────────────────────────────────────────────────┘
```

## Available Gadgets

Currently configured gadgets:
1. **Trace Exec** (trace_exec) - Process execution tracing
2. **Trace TCP** (trace_tcp) - TCP connection monitoring
3. **Snapshot Process** (snapshot_process) - Process snapshot

## Next Steps

### Step 1: Install Dependencies

```bash
cd frontend
npm install
```

Or use the Makefile:
```bash
make install-deps-frontend
```

### Step 2: Test Locally

#### Option A: Development Mode (Recommended)

1. Start the backend:
```bash
make dev-backend
# Or: cd backend && go run cmd/server/main.go
```

2. In a new terminal, start the frontend:
```bash
make dev-frontend
# Or: cd frontend && npm run dev
```

3. Open your browser to `http://localhost:3000`

#### Option B: Build and Deploy to Kubernetes

1. Build the Docker images:
```bash
./build.sh
# Or: make build
```

2. Import to k3s (if using k3s):
```bash
make import-k3s
```

3. Deploy:
```bash
./deploy.sh
# Or: make deploy
```

4. Port forward:
```bash
make port-forward
```

5. Open your browser to `http://localhost:3000`

## Features Implemented

### ✅ Catalog View
- Card-based gadget browser
- Category filtering
- Visual indicators for running gadgets
- Hover effects and animations

### ✅ Runner Interface
- Configuration panel with namespace and pod filters
- TCP-specific options (accept-only, connect-only, failure-only)
- Start/Stop controls
- Active configuration display

### ✅ Global State Management
- Active Operations sidebar panel
- Running gadget tracking
- Smart resume when selecting active gadgets
- Visual badges and pulse animations

### ✅ Backend Integration
- Uses existing REST API for session management
- WebSocket streaming for real-time data
- Session tracking and management
- Error handling

### ✅ Output Display
- Visual and Raw JSON tabs
- Reuses existing GadgetOutput component
- Export JSON functionality

## Testing Checklist

When testing the new UI, verify:

- [ ] Catalog displays all 3 gadgets
- [ ] Category filters work (All, Trace, Snapshot)
- [ ] Clicking a gadget opens the runner
- [ ] Configuration panel accepts input
- [ ] Start button creates a session
- [ ] WebSocket connects and displays output
- [ ] Running badge appears in catalog
- [ ] Active Operations panel shows running gadgets
- [ ] Clicking running gadget in sidebar switches to it
- [ ] Stop button terminates the session
- [ ] Export JSON downloads the data
- [ ] Visual/Raw JSON tabs switch views
- [ ] Error messages display correctly

## Trace TCP Specific Features

When testing `trace_tcp`:
1. Open the Trace TCP gadget
2. Configure:
   - Namespace: `default` (or your test namespace)
   - Pod Name: Leave empty for all pods, or specify one
   - TCP Filters:
     - ☑️ **Accept Only**: Show only accepted connections
     - ☑️ **Connect Only**: Show only outgoing connections
     - ☑️ **Failure Only**: Show only failed connections
3. Click "Start Gadget"
4. Watch real-time TCP events appear
5. Use the existing visualization (Flow Diagram, Connection Lifecycle)

## File Changes Summary

### Modified Files
- `frontend/package.json` - Added dependencies
- `frontend/src/index.css` - Added Tailwind directives
- `frontend/src/App.tsx` - Complete rewrite with catalog layout

### New Files
- `frontend/tailwind.config.js`
- `frontend/postcss.config.js`
- `frontend/src/components/GadgetCard.tsx`
- `frontend/src/components/Runner.tsx`

### Unchanged (Backend)
- No backend changes required
- All existing APIs work as-is
- `backend/` directory untouched

## Troubleshooting

### Issue: Tailwind classes not working
**Solution**: Make sure you ran `npm install` to install Tailwind and its dependencies.

### Issue: Icons not showing
**Solution**: Verify `lucide-react` was installed: `npm list lucide-react`

### Issue: WebSocket connection fails
**Solution**:
- Check backend is running on port 8080
- Verify `frontend/vite.config.ts` proxy settings
- Check browser console for errors

### Issue: No gadgets appearing
**Solution**: Check browser console for React errors. Make sure all imports are correct.

### Issue: Backend connection refused
**Solution**:
- Ensure backend is running: `make dev-backend`
- Check backend logs for errors
- Verify kubectl-gadget is installed in the backend container

## Next Features to Implement

Based on your notes in `new-ui-ideas/notes.txt`, here are suggested next steps:

1. **Enhanced Visualizations for trace_tcp**:
   - Bar charts showing bandwidth by pod (using Recharts)
   - Network topology diagram
   - Connection states table

2. **Add More Gadgets**:
   - trace_dns
   - top_tcp (with bar chart visualization)
   - top_file
   - profile_cpu (with flame graphs)
   - audit_seccomp

3. **Problem-First Search**:
   - Add search bar to catalog
   - Implement keyword → gadget suggestions

4. **Enhanced Filtering**:
   - Service IP dropdowns
   - Node selector with cluster map

5. **Advanced Features**:
   - Pause/Resume streaming
   - Time range filters
   - Export to CSV/PDF
   - Session history

## Resources

- **Inspektor Gadget Docs**: https://github.com/inspektor-gadget/inspektor-gadget
- **Tailwind CSS**: https://tailwindcss.com/docs
- **Recharts**: https://recharts.org/
- **Lucide Icons**: https://lucide.dev/icons/

## Support

If you encounter any issues:
1. Check the browser console for errors
2. Check backend logs: `make logs-backend`
3. Verify dependencies are installed: `npm list`
4. Ensure kubectl-gadget is accessible in the backend

---

**Status**: ✅ Integration Complete - Ready for Testing

The new UI is fully integrated and ready to be tested. Run `npm install` in the frontend directory and start testing!
