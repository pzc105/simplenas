import React, { useState } from 'react';
import { Typography, Drawer, IconButton } from '@mui/material';
import { styled } from "@mui/material/styles";

const SideUtilsRoot = styled('div')(({ theme }) => ({
  display: 'flex',
}))
const SideIcon = styled(IconButton)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  borderRadius: '50%',
  width: '64px',
  height: '32px',
  padding: '0',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
}))
const SideDrawer = styled(Drawer)(({ theme }) => ({
  width: 240,
  flexShrink: 0,
}))

export default function SideUtils({ name, child }) {
  const [open, setOpen] = useState(false);
  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  return (
    <SideUtilsRoot >
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-start' }}>
        <SideIcon onClick={handleDrawerToggle}>
          <Typography gutterBottom variant="body1" component="div" noWrap>
            {name}
          </Typography>
        </SideIcon>
      </div>

      <SideDrawer
        open={open}
        onClose={handleDrawerToggle}>
        {child}
      </SideDrawer>
    </SideUtilsRoot>
  )
}