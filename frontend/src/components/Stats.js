import React, { useState, useEffect, useCallback, useMemo } from 'react';
import axios from 'axios';
import { 
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
  CircularProgress,
  RadioGroup,
  FormControlLabel,
  Radio
} from '@mui/material';
import { styled } from '@mui/material/styles';
import FileDownloadIcon from '@mui/icons-material/FileDownload';

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
  const [filterOption, setFilterOption] = useState('all');
  const [timeRange, setTimeRange] = useState('1h');
  const [rawStats, setRawStats] = useState({ headers: [], data: [] });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadStats = useCallback(async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/stats', {
        params: { timeRange }
      });
      setRawStats(response.data);
      setError(null);
    } catch (error) {
      console.error("Error fetching stats:", error);
      setError('Error fetching stats');
    } finally {
      setLoading(false);
    }
  }, [timeRange]);

  useEffect(() => {
    loadStats();
  }, [loadStats, timeRange]);

  const filteredStats = useMemo(() => {
    const trafficIndex = rawStats.headers.findIndex(header => header.toLowerCase().includes('traffic') || header.toLowerCase().includes('count'));
    
    if (trafficIndex === -1) {
      return rawStats.data;
    }

    return rawStats.data.filter(row => {
      const trafficValue = row[trafficIndex];
      switch (filterOption) {
        case 'withTraffic':
          return trafficValue !== '0' && trafficValue !== '-';
        case 'withoutTraffic':
          return trafficValue === '0' || trafficValue === '-';
        default:
          return true;
      }
    });
  }, [rawStats, filterOption]);

  const handleFilterChange = (event) => {
    setFilterOption(event.target.value);
  };

  const handleTimeRangeChange = (event) => {
    setTimeRange(event.target.value);
  };

  const handleDownloadCSV = () => {
    const csvContent = [
      rawStats.headers.join(','),
      ...filteredStats.map(row => row.join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `apigee_stats_${timeRange}_${filterOption}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
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
          <Grid container spacing={2} alignItems="flex-end">
            <Grid item xs={12} sm={6}>
              <FormControl component="fieldset">
                <RadioGroup
                  row
                  value={filterOption}
                  onChange={handleFilterChange}
                >
                  <FormControlLabel value="all" control={<Radio />} label="All" />
                  <FormControlLabel value="withTraffic" control={<Radio />} label="With Traffic" />
                  <FormControlLabel value="withoutTraffic" control={<Radio />} label="Without Traffic" />
                </RadioGroup>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={4}>
              <FormControl fullWidth>
                <InputLabel>Time Range</InputLabel>
                <Select
                  value={timeRange}
                  onChange={handleTimeRangeChange}
                >
                  <MenuItem value="1h">1 hour</MenuItem>
                  <MenuItem value="6h">6 hours</MenuItem>
                  <MenuItem value="12h">12 hours</MenuItem>
                  <MenuItem value="1d">1 day</MenuItem>
                  <MenuItem value="7d">7 days</MenuItem>
                  <MenuItem value="14d">14 days</MenuItem>
                  <MenuItem value="30d">30 days</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={2}>
              <Button 
                variant="contained" 
                color="secondary" 
                fullWidth
                onClick={handleDownloadCSV}
                startIcon={<FileDownloadIcon />}
                disabled={loading || filteredStats.length === 0}
              >
                Download CSV
              </Button>
            </Grid>
          </Grid>
        </Grid>
        <Grid item>
          {loading ? (
            <CircularProgress />
          ) : error ? (
            <Typography color="error">{error}</Typography>
          ) : (
            <TableContainer component={Paper} elevation={3}>
              <Table>
                <TableHead>
                  <TableRow>
                    {rawStats.headers.map((header, index) => (
                      <StyledTableCell key={index}>{snakeToTitleCase(header)}</StyledTableCell>
                    ))}
                  </TableRow>
                </TableHead>
                <TableBody>
                  {filteredStats.map((row, rowIndex) => (
                    <TableRow key={rowIndex} hover>
                      {row.map((cell, cellIndex) => (
                        <TableCell key={cellIndex}>{cell}</TableCell>
                      ))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Grid>
      </Grid>
    </Container>
  );
}

export default Stats;