import React, { useState } from 'react';
import Pagination from '@mui/material/Pagination';
import Stack from '@mui/material/Stack';

export default function UnifiedPage({ PageTotalCount, PageNum, onPage }) {
  return (
    <Stack spacing={2}>
      {/* 分页组件 */}
      <Pagination
        count={PageTotalCount} // 总页数
        page={PageNum}
        onChange={(event, value) => onPage(value)}
        showFirstButton
        showLastButton
      />
    </Stack>
  );
}
