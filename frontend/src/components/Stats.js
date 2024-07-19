import React, { useState } from 'react';
import axios from 'axios';
import { 
  TextField, 
  Select, 
  MenuItem, 
  Button, 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Paper,
  Grid,
  Typography,
  FormControl,
  InputLabel,
  Container,
  CircularProgress
} from '@mui/material';
import { styled } from '@mui/material/styles';

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.common.white,
  fontWeight: 'bold',
}));

function snakeToTitleCase(str) {
  return str.split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

function Stats() {
  const [proxyName, setProxyName] = useState('');
  const [timeRange, setTimeRange] = useState('1h');
  const [stats, setStats] = useState({ headers: [], data: [] });
  const [tableLoading, setTableLoading] = useState(false);
  const [error, setError] = useState(null);

  const loadStats = async () => {
    setTableLoading(true);
    setError(null);
    try {
      const response = await axios.get('/api/stats', {
        params: { proxyName, timeRange }
      });
      setStats(response.data);
    } catch (error) {
      console.error("Error fetching stats:", error);
      setError('Error fetching stats');
    } finally {
      setTableLoading(false);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    loadStats();
  };

  return (
    <Container maxWidth="lg">
      <Grid container direction="column" spacing={3}>
        <Grid item>
          <Typography variant="h4" gutterBottom align="center">
            Apigee Proxy Stats (max duration)
          </Typography>
        </Grid>
        <Grid item>
          <form onSubmit={handleSubmit}>
            <Grid container spacing={2} alignItems="flex-end">
              <Grid item xs={12} sm={4}>
                <TextField
                  fullWidth
                  label="Proxy Name (optional)"
                  value={proxyName}
                  onChange={(e) => setProxyName(e.target.value)}
                />
              </Grid>
              <Grid item xs={12} sm={4}>
                <FormControl fullWidth>
                  <InputLabel>Time Range</InputLabel>
                  <Select
                    value={timeRange}
                    onChange={(e) => setTimeRange(e.target.value)}
                  >
                    <MenuItem value="1h">1 hour</MenuItem>
                    <MenuItem value="6h">6 hours</MenuItem>
                    <MenuItem value="12h">12 hours</MenuItem>
                    <MenuItem value="1d">1 day</MenuItem>
                    <MenuItem value="7d">7 days</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} sm={4}>
                <Button 
                  type="submit" 
                  variant="contained" 
                  color="primary" 
                  fullWidth
                >
                  Search
                </Button>
              </Grid>
            </Grid>
          </form>
        </Grid>
        
        {error && (
          <Grid item>
            <Typography color="error">{error}</Typography>
          </Grid>
        )}
        
        <Grid item>
          <TableContainer component={Paper} elevation={3} style={{ position: 'relative' }}>
            {tableLoading && (
              <div style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                background: 'rgba(255, 255, 255, 0.7)',
                zIndex: 1,
              }}>
                <CircularProgress />
              </div>
            )}
            <Table>
              <TableHead>
                <TableRow>
                  {stats.headers.map((header, index) => (
                    <StyledTableCell key={index}>{snakeToTitleCase(header)}</StyledTableCell>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {stats.data.map((row, rowIndex) => (
                  <TableRow key={rowIndex} hover>
                    {row.map((cell, cellIndex) => (
                      <TableCell key={cellIndex}>{cell}</TableCell>
                    ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Grid>
      </Grid>
    </Container>
  );
}

export default Stats;