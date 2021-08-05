import React from 'react';
import PropTypes from 'prop-types';
import { Popover, } from '@material-ui/core';
import { ButtonScore, PopoverText, } from './styles/DashboardStyles';

// Represents the reputation of a user
const ReputationElement = (props) => {
  const { repScore, } = props;
  const [anchorEl, setAnchorEl] = React.useState(null);
  console.log(`Score --> ${repScore}`);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const starterString = `Carbon Reputation: carbon emissions of seller at offer 
  creation.`;

  const open = Boolean(anchorEl);
  const id = open ? 'simple-popover' : undefined;

  return (
    <div>
      <ButtonScore
        aria-describedby={id}
        variant="contained"
        color="primary"
        onClick={handleClick}
        score={repScore}
      >
        Reputation: {repScore}
      </ButtonScore>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <PopoverText>
          {repScore >= 1000
            ? starterString +
            '\nWarning: bad carbon reputation.'
            : ''}
          {repScore < 1000 && repScore >= 500
            ? starterString +
            '\nAverage carbon reputation.'
            : ''}
          {repScore < 500 ? starterString + '\nGreat carbon reputation.' : ''}
        </PopoverText>
      </Popover>
    </div>
  );
};

export default ReputationElement;

ReputationElement.propTypes = {
  repScore: PropTypes.number.isRequired,
};
