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

function Deployments() {
  const [deployments, setDeployments] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    axios.get('/api/deployments')
      .then(response => {
        setDeployments(response.data);
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