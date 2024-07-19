import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { 
  Table, 
  TableBody, 
  TableCell, 
  Typography, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Paper, 
  CircularProgress, 
  Grid,
  Container
} from '@mui/material';
import { styled } from '@mui/material/styles';

// Styled components for better header formatting
const StyledTableCell = styled(TableCell)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.common.white,
  fontWeight: 'bold',
}));

const environmentOrder = ['dev', 'test', 'uat', 'prod', 'wc-prod'];

function reorderDeploymentData(headers, data) {
  const nameIndex = headers.findIndex(h => h.toLowerCase() === 'name');
  const reorderedHeaders = ['Name', ...environmentOrder];
  
  const reorderedData = data.map(row => {
    const newRow = [row[nameIndex]];
    environmentOrder.forEach(env => {
      const envIndex = headers.findIndex(h => h.toLowerCase() === env);
      newRow.push(envIndex !== -1 ? row[envIndex] : '-');
    });
    return newRow;
  });

  return { headers: reorderedHeaders, data: reorderedData };
}

function Deployments() {
  const [deployments, setDeployments] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    axios.get('/api/deployments')
      .then(response => {
        const reorderedDeployments = reorderDeploymentData(response.data.headers, response.data.data);
        setDeployments(reorderedDeployments);
        setLoading(false);
      })
      .catch(error => {
        setError('Error fetching deployments');
        setLoading(false);
      });
  }, []);

  if (loading) return (
    <Container>
      <Grid container justifyContent="center" alignItems="center" style={{ minHeight: '100vh' }}>
        <CircularProgress />
      </Grid>
    </Container>
  );
  
  if (error) return (
    <Container>
      <Grid container justifyContent="center" alignItems="center" style={{ minHeight: '100vh' }}>
        <Typography color="error">{error}</Typography>
      </Grid>
    </Container>
  );

  return (
    <Container maxWidth="lg">
      <Grid container direction="column" spacing={3}>
        <Grid item>
          <Typography variant="h4" gutterBottom align="center">
            Apigee Proxy Deployments
          </Typography>
        </Grid>
        <Grid item>
          <TableContainer component={Paper} elevation={3}>
            <Table>
              <TableHead>
                <TableRow>
                  {deployments.headers.map((header, index) => (
                    <StyledTableCell key={index}>{header}</StyledTableCell>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {deployments.data.map((row, rowIndex) => (
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

export default Deployments;