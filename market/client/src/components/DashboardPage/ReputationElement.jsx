import React from 'react';
import PropTypes from 'prop-types';
import { Popover, } from '@material-ui/core';
import { ButtonScore, PopoverText, } from './styles/DashboardStyles';
import { useHistory, } from 'react-router-dom';
import Link from '@mui/material/Link';

// Represents the reputation of a user
const ReputationElement = (props) => {
  const { repScore, username, } = props;
  const [anchorEl, setAnchorEl] = React.useState(null);
  console.log(`Score --> ${repScore}`);
  console.log(`Producer --> ${username}`);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const starterString = `Seller Reputation: score provided to the user in relation 
  to Environmental, Social and Governance contributions.`;

  const open = Boolean(anchorEl);
  const id = open ? 'simple-popover' : undefined;
  const history = useHistory();
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
          {starterString + '\n'}
          <b>
            {repScore <= -500
              ? '\nWarning: bad reputation.'
              : ''}
            {repScore < -250 && repScore >= -500
              ? '\nAverage reputation.'
              : ''}
            {repScore > -250 ? '\nGreat reputation.' : ''}
          </b>
          <Link onClick={() => {
            history.push('/production',
              { name: username, });
          }}> Visit Reputation Breakdown of {username} </Link>
        </PopoverText>
      </Popover>
    </div>
  );
};

export default ReputationElement;

ReputationElement.propTypes = {
  repScore: PropTypes.number.isRequired,
  username: PropTypes.string.isRequired,
};
