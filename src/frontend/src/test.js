import React, { useState } from 'react';
import { Button, Popover, Box } from '@mui/material';
import Draggable from 'react-draggable';
import { Container } from '@mui/system';
import CloseIcon from '@mui/icons-material/Close';

export function Test() {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'floating-window' : undefined;

  return (
    <Container>
      <Button variant="contained" color="primary" onClick={handleClick}>
        Open Floating Window
      </Button>
      <Draggable>
        <Popover
          id={id}
          open={open}
          anchorEl={anchorEl}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
        >
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', pr: 1 }}>
            <Button size="small" color="secondary" onClick={handleClose}>
              <CloseIcon />
            </Button>
          </Box>
          <div style={{ padding: '20px' }}>
            <p>This is the content of the floating window.</p>
          </div>
        </Popover>
      </Draggable>

    </Container>
  );
}
