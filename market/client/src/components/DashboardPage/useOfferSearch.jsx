import { useEffect, useState, } from 'react';

/**
 * Component for searching offers
 * @param {*} token paginationToken
 * @param {*} authToken for authorisation
 * @returns loading, err, offers, hasMore, paginationToken
 */
const useOfferSearch = (token, authToken, apiFun) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [offers, setOffers] = useState([]);
  const [hasMore, setHasMore] = useState(false);
  const [paginationToken, setPaginationToken] = useState(token);

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;
    const offerRetrieval = async () => {
      setLoading(true);
      setError('');
      try {
        const response = await apiFun(authToken, signal, token);
        setOffers((prev) => {
          return [...new Set([...prev, ...response.records])];
        });
        setHasMore(response.fetchedRecordsCount >= 10);
        setPaginationToken(response.bookmark);
        setLoading(false);
      } catch (err) {
        if (!(err instanceof DOMException)) {
          setError(err.message);
        }
      }
    };
    offerRetrieval();
    return () => controller.abort();
  }, [token]);
  return { loading, error, offers, hasMore, paginationToken, };
};

export default useOfferSearch;
