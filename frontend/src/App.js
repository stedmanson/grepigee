import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import Home from './components/Home';
import Stats from './components/Stats';
import Deployments from './components/Deployments';
import Grep from './components/Grep';

function App() {
  return (
    <Router>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            Grepigee
          </Typography>
          <Button color="inherit" component={Link} to="/">Home</Button>
          <Button color="inherit" component={Link} to="/stats">Stats</Button>
          <Button color="inherit" component={Link} to="/deployments">Deployments</Button>
          <Button color="inherit" component={Link} to="/grep">Grep</Button>
        </Toolbar>
      </AppBar>

      <Box mt={4}> {/* This adds space below the AppBar */}
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/stats" element={<Stats />} />
          <Route path="/deployments" element={<Deployments />} />
          <Route path="/grep" element={<Grep />} />
        </Routes>
      </Box>
    </Router>
  );
}

export default App;