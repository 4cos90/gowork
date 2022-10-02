import service from './index'

export const listBeer = (pageNum, pageSize) => {
    return service.get("/v1/catalog/beers", {
        pageNum, pageSize
    })
};

export const getBeerDetail = (id) => {
    return service.get("/v1/catalog/beers/"+id)
};
