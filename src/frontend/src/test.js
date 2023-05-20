import * as React from 'react';
import { List, ListItemButton } from '@mui/material';

export function MultiSelectList() {
  const [selectedIndex, setSelectedIndex] = React.useState(null);

  const handleListItemMouseDown = (event, index) => {
    if (event.ctrlKey) {
      // 如果按下了 Ctrl 键，则切换索引的选中状态
      setSelectedIndex((prevSelectedIndex) => {
        if (prevSelectedIndex && prevSelectedIndex.includes(index)) {
          return prevSelectedIndex.filter((i) => i !== index);
        } else {
          return prevSelectedIndex ? [...prevSelectedIndex, index] : [index];
        }
      });
    } else {
      // 如果没有按下 Ctrl 键，则只选中当前索引
      setSelectedIndex([index]);
    }
  };

  const handleListItemMouseUp = () => {
    // 在鼠标抬起时重置选中状态
    setSelectedIndex(null);
  };

  return (
    <List>
      {[0, 1, 2, 3, 4].map((index) => (
        <ListItemButton
          key={index}
          selected={selectedIndex && selectedIndex.includes(index)}
          onMouseDown={(event) => handleListItemMouseDown(event, index)}
          onMouseUp={handleListItemMouseUp}
        >
          Item {index + 1}
        </ListItemButton>
      ))}
    </List>
  );
}
