var app = angular.module('speedrunRouting', ['ngVis', 'ngAnimate', 'ui.bootstrap']);

app.controller('RouteCtrl', function($log, $http, VisDataSet, $location) {

    var DEBUG = true;

    var g = this;

    g.rewards = [];
    g.nodes = []
    g.edges = [];
    g.startNode = "";
    g.endNode = "";
    g.password = "";
    g.shortestPath = [];

    g.rewardBeingEdited = undefined;
    g.nodeBeingEdited = undefined;
    g.edgeBeingEdited = undefined;

    /* vis network data start */
    var networkNodes = new VisDataSet();
    var networkEdges = new VisDataSet();
    g.network = {
        nodes: networkNodes,
        edges: networkEdges
    };
    g.options = {
        height: '100%'
    };
    g.events = {};
    /* vis network data end */

    g.resetReward = function() {
        g.rewardBeingEdited = undefined;
        g.reward = {
            edit: true,
            id: "",
            unique: false,
            isA: "",
            errors: []
        }
    };

    resetRewardRef = function(obj) {
        obj.rewardRef = {
            rewardId: "",
            quantity: ""
        }
    };

    g.resetNode = function() {
        g.nodeBeingEdited = undefined;
        g.node = {
            edit: true,
            id: "",
            revisitable: false,
            rewards: [],
            errors: []
        }
        resetRewardRef(g.node);
    };

    resetEdgeWeight = function() {
        g.edge.weight = {
            description: ""
        }
        resetRewardRef(g.edge.weight);
    };

    g.resetEdge = function() {
        g.edgeBeingEdited = undefined;
        g.edge = {
            edit: true,
            from: "",
            to: "",
            weights: [],
            errors: []
        }
        resetEdgeWeight();
    };

    containsEdge = function(list, from, to) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].from === from && list[i].to === to) {
                return true;
            }
        }
        return false;
    };

    contains = function(list, id) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].id === id) {
                return true;
            }
        }
        return false;
    };

    containsReward = function(list, id) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].rewardId === id) {
                return true;
            }
        }
        return false;
    };

    removeObj = function(list, obj) {
        var index = list.indexOf(obj);
        if(index > -1) {
            list.splice(index, 1);
        }
    };

    g.canRewardBeRemoved = function(id) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            for(var j = 0; j < edge.weights.length; j++) {
                var weight = edge.weights[j];
                for(var k = 0; k < weight.requirements.length; k++) {
                    if(weight.requirements[k].rewardId === id) {
                        return false;
                    }
                }
            }
        }
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            for(var j = 0; j < node.rewards.length; j++) {
                if(node.rewards[j].rewardId === id) {
                    return false;
                }
            }
        }
        for(var i = 0; i < g.rewards.length; i++) {
            if(g.rewards[i].isA === id) {
                return false;
            }
        }
        return true;
    };

    g.canNodeBeRemoved = function(id) {
        for(var i = 0; i < g.edges.length; i++) {
            if(g.edges[i].from === id || g.edges[i].to === id) {
                return false;
            }
        }
        return true;
    };

    toggleEdit = function(reward, node, edge) {
        g.reward.edit = reward;
        g.node.edit = node;
        g.edge.edit = edge;
    };

    log = function(msg) {
        if(DEBUG) {
            console.log(msg);
        }
    };

    updateNodeRewardReferences = function(oldId, newId) {
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            for(var j = 0; j < node.rewards.length; j++) {
                var nodeRew = node.rewards[j];
                if(nodeRew.rewardId === oldId) {
                    nodeRew.rewardId = newId;
                    log("updated node reward reference");
                }
            }
        }
    };

    updateEdgeRequirementReferences = function(oldId, newId) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            for(var j = 0; j < edge.weights.length; j++) {
                var weight = edge.weights[j];
                for(var k = 0; k < weight.requirements.length; k++) {
                    var weightReq = weight.requirements[k];
                    if(weightReq.rewardId === oldId) {
                        weightReq.rewardId = newId;
                        log("updated edge weight requirement reference");
                    }
                }
            }
        }
    };

    updateEdgeNodesReferences = function(oldId, newId) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            if(edge.from === oldId) {
                edge.from = newId;
                log("updated edge from node reference");
            }
            if(edge.to === oldId) {
                edge.to = newId;
                log("updated edge to node reference");
            }
        }
    };

    getNodeIndex = function(id) {
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            if(node.id === id) {
                return i;
            }
        }
        return -1;
    };

    getEdgeIndex = function(from, to) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            if(edge.from === from && edge.to === to) {
                return i;
            }
        }
        return -1;
    };

    g.len = function(list) {
        if (list) {
            return list.length;
        }
        return 0;
    };

    g.addReward = function() {
        g.reward.errors = [];
        var error = false;
        var id = g.reward.id;
        var isA = g.reward.isA;
        if(!id) {
            g.reward.errors.push("The reward name is not set.");
            error = true;
        }
        if(isA && !contains(g.rewards, isA)){
            g.reward.errors.push(isA + " is not a valid reward reference.");
            error = true;
        }

        if(g.rewardBeingEdited) {
            log("EDIT MODE");
            var oldId = g.rewardBeingEdited.id;
            var hasChanged = oldId !== id;
            if(hasChanged && contains(g.rewards, id)) {
                g.reward.errors.push("The updated reward name already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            g.rewardBeingEdited.unique = g.reward.unique;
            g.rewardBeingEdited.isA = isA;
            if(hasChanged) {
                g.rewardBeingEdited.id = id;
                updateNodeRewardReferences(oldId, id);
                updateEdgeRequirementReferences(oldId, id);
            }
            g.resetReward();
        } else {
            log("ADD MODE");
            if(contains(g.rewards, id)){
                g.reward.errors.push("The reward name already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            g.rewards.push({
                id: id,
                unique: g.reward.unique,
                isA: isA
            });
            g.resetReward();
        }
    };

    g.removeNodeReward = function(index) {
        g.node.rewards.splice(index, 1);
    };

    g.removeEdgeWeight = function(index) {
        g.edge.weights.splice(index, 1);
    };

    g.removeEdgeWeightRequirement = function(wIndex, index) {
        g.edge.weights[wIndex].requirements.splice(index, 1);
    };

    g.addNodeReward = function() {
        g.node.errors = [];
        var error = false;
        var id = g.node.rewardRef.rewardId;

        if(!id) {
            g.node.errors.push("The reward name cannot be empty.");
            error = true;
        }
        if(id && !contains(g.rewards, id)) {
            g.node.errors.push(id + " is not a valid reward reference.");
            error = true;
        }
        if(id && containsReward(g.node.rewards, id)) {
            g.node.errors.push(id + " is already defined in the list.");
            error = true;
        }
        if(error) {
            return;
        }

        g.node.rewards.push({
            rewardId: id,
            quantity: (parseInt(g.node.rewardRef.quantity) || 1)//check if this can be removed if not set
        });
        resetRewardRef(g.node);
    };

    g.addEdgeWeightRequirement = function(wIndex) {
        var id = g.edge.weight.rewardRef.rewardId;
        if(id && !containsReward(g.edge.weights[wIndex].requirements, id) && contains(g.rewards, id)) {
            g.edge.weights[wIndex].requirements.push({
                rewardId: id,
                quantity: (parseInt(g.edge.weight.rewardRef.quantity) || 1)//check if this can be removed if not set
            });
            resetRewardRef(g.edge.weight);
        } else {
            log("Error, containsReward or contains something something")
        }
    };

    g.addEdgeWeight = function() {
        var weight = g.edge.weight;
        g.edge.weights.push({
            requirements: [],
            description: weight.description,
            time: ""
        });
        resetEdgeWeight();
    };

    g.addNode = function() {
        g.node.errors = [];
        var error = false;
        var id = g.node.id;

        if(!id) {
            g.node.errors.push("The node name cannot be empty.");
            error = true;
        }
        if(g.node.rewardRef.rewardId !== ""){
            g.node.errors.push("There cannot be an unfinished reward.");
            error = true;
        }

        if(g.nodeBeingEdited) {
            log("EDIT MODE");
            var oldId = g.nodeBeingEdited.id;
            var isNewIdOk = oldId !== id && !contains(g.nodes, id);

            if(!(oldId === id || isNewIdOk)) {
                g.node.errors.push("The updated node name already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            g.nodeBeingEdited.revisitable = g.node.revisitable;
            g.nodeBeingEdited.rewards = g.node.rewards;
            if(isNewIdOk) {
                log("newId ok");
                g.nodeBeingEdited.id = id;
                updateEdgeNodesReferences(oldId, id);
                networkNodes.update({
                    id: g.nodes.indexOf(g.nodeBeingEdited),
                    label: id,
                    color: 'lightblue'
                });
            }
            g.resetNode();
        } else {
            log("ADD MODE");

            if(contains(g.nodes, id)) {
                g.node.errors.push("The node name already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            g.nodes.push({
                id: id,
                revisitable: g.node.revisitable,
                rewards: g.node.rewards
            });
            networkNodes.add({
                id: g.nodes.length-1,
                label: id,
                color: 'lightblue'
            });
            g.resetNode();
        }
    };

    g.addEdge = function() {
        g.edge.errors = [];
        var error = false;
        var from = g.edge.from;
        var to = g.edge.to;

        if(!from) {
            g.edge.errors.push("From node cannot be empty.");
            error = true;
        }
        if(!to) {
            g.edge.errors.push("To node cannot be empty.");
            error = true;
        }
        if(!contains(g.nodes, from)) {
            g.edge.errors.push("From node " + from + " doesn't exist.");
            error = true;
        }
        if(!contains(g.nodes, to)) {
            g.edge.errors.push("To node " + to + " doesn't exist.");
            error = true;
        }
        if(g.edge.weight.description) {
            g.edge.errors.push("There cannot be an unfinished weight.");
            error = true;
        }
        for(var i = 0; i < g.edge.weights.length; i++) {
            var weight = g.edge.weights[i];
            if(!/^\d*[:]{0,1}\d*[:]{0,1}\d*[.]{0,1}\d*$/.test(weight.time)) {
                g.edge.errors.push("The time of the weight " + weight.description + " is not on the correct format.");
                error = true;
            }
        }

        if(g.edgeBeingEdited) {
            log("EDIT MODE");
            var oldFrom = g.edgeBeingEdited.from;
            var oldTo = g.edgeBeingEdited.to;
            var hasChanged = (oldFrom !== from) || (oldTo !== to);

            if(hasChanged && containsEdge(g.edges, from, to)) {
                g.edge.errors.push("The edge from " + from + " to " + to + " already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            g.edgeBeingEdited.weights = g.edge.weights;
            if(hasChanged) {
                g.edgeBeingEdited.from = from;
                g.edgeBeingEdited.to = to;
                networkEdges.update({
                    id: g.edges.indexOf(g.edgeBeingEdited),
                    from: getNodeIndex(g.edgeBeingEdited.from),
                    to: getNodeIndex(g.edgeBeingEdited.to),
                    arrows: 'to',
                    color: 'lightblue'
                });
            }
            g.resetEdge();
        } else {
            log("ADD MODE");

            if(containsEdge(g.edges, from, to)) {
                g.edge.errors.push("The edge from " + from + " to " + to + " already exists.");
                error = true;
            }
            if(error) {
                return;
            }

            var edge = {
                from: from,
                to: to,
                weights: g.edge.weights
            };
            g.edges.push(edge);
            networkEdges.add({
                id: g.edges.length-1,
                from: getNodeIndex(edge.from),
                to: getNodeIndex(edge.to),
                arrows: 'to',
                color: 'lightblue'
            });
            g.resetEdge();
        }
    };

    g.edgeRowSpan = function(weights) {
        var span = 0;
        for (var i = 0; i < g.len(weights); i++) {
            for (var j = 1; j < g.len(weights[i].requirements); j++) {
                span++;
            }
            span++;
        }
        return span;
    };

    g.removeReward = function(reward) {
        if(g.canRewardBeRemoved(reward.id)) {
            removeObj(g.rewards, reward);
        }
    };

    g.removeNode = function(node) {
        if(g.canNodeBeRemoved(node.id)) {
            networkNodes.remove({
                id: g.nodes.indexOf(node)
            });
            removeObj(g.nodes, node);
        } else {
            log("Nonono");
        }
    };

    g.removeEdge = function(edge) {
        networkEdges.remove({
            id: g.edges.indexOf(edge)
        });
        removeObj(g.edges, edge);
    };

    g.editReward = function(reward) {
        g.resetReward();
        toggleEdit(true, false, false);
        g.rewardBeingEdited = reward;
        g.reward.id = reward.id;
        g.reward.unique = reward.unique;
        g.reward.isA = reward.isA;
    };

    g.editNode = function(node) {
        g.resetNode();
        toggleEdit(false, true, false);
        g.nodeBeingEdited = node;
        g.node.id = node.id;
        g.node.revisitable = node.revisitable;
        g.node.rewards = angular.copy(node.rewards);
    };

    g.editEdge = function(edge) {
        g.resetEdge();
        toggleEdit(false, false, true);
        g.edgeBeingEdited = edge;
        g.edge.from = edge.from;
        g.edge.to = edge.to;
        g.edge.weights = angular.copy(edge.weights);
    };

    resetColors = function() {
        for(var i = 0; i < networkNodes.length; i++) {
            networkNodes.update({
                id: i,
                color: 'lightblue'
            });
        }
        for(var i = 0; i < networkEdges.length; i++) {
            networkEdges.update({
                id: i,
                color: 'lightblue'
            });
        }
    }

    g.saveGraph = function() {
        resetColors();
        g.shortestPath = [];
        $http({
            method: 'POST',
            url: '/graph/' + g.name + '/' + g.password,
            data: {
                rewards: g.rewards,
                nodes: g.nodes,
                edges: g.edges,
                startId: g.startNode,
                endId: g.endNode
            }
        }).then(function successCallback(response) {
            var data = response.data;
            if(!data) {
                return;
            }
            g.shortestPath = data;
            for(var i = 0; i < data.length; i++) {
                if(i != 0) {
                    networkEdges.update({
                        id: getEdgeIndex(data[i-1], data[i]),
                        color: 'lime'
                    });
                }
                networkNodes.update({
                    id: getNodeIndex(data[i]),
                    color: 'lime'
                });
            }
        }, function errorCallback(response) {
            log("404 I guess...");
        });
    }

    loadGraph = function() {
        $http({
            method: 'GET',
            url: '/graph/' + g.name
        }).then(function successCallback(response) {
            var data = response.data;
            if(data.rewards) {
                g.rewards = data.rewards;
            }
            if(data.nodes) {
                g.nodes = data.nodes;
            }
            if(data.edges) {
                g.edges = data.edges;
            }
            if(data.startId) {
                g.startNode = data.startId;
            }
            if(data.endId) {
                 g.endNode = data.endId;
            }
            for(var i = 0; i < g.nodes.length; i++) {
                networkNodes.add({
                    id: i,
                    label: g.nodes[i].id,
                    color: 'lightblue'
                });
            }
            for(var i = 0; i < g.edges.length; i++) {
                networkEdges.add({
                    id: i,
                    from: getNodeIndex(g.edges[i].from),
                    to: getNodeIndex(g.edges[i].to),
                    arrows: 'to',
                    color: 'lightblue'
                });
            }
        }, function errorCallback(response) {
            log("404 I guess...");
        });
    }

    init = function() {
        var url = window.location.href;
        var parameter = 'g';
        var regex = new RegExp("[?&]" + parameter + "(=([^&#]*)|&|#|$)", "i");
        var results = regex.exec(url);
        if (!results || !results[2]) {
            log("Couldn't find key");
            return;
        }
        g.name = decodeURIComponent(results[2].replace(/\+/g, " "));
        loadGraph();
        g.reward.edit = false;
        g.node.edit = false;
        g.edge.edit = false;
    }

    g.resetReward();
    g.resetNode();
    g.resetEdge();
    init();
});

app.controller('CreateCtrl', function($http) {

    var c = this;

    c.name = "";
    c.password = "";
    c.livesplit = "";
    c.createError = 0;

    c.createGraph = function() {
        c.createError = 0;
        if(!c.name || !/^[A-Za-z0-9-_]*$/.test(c.name)) {
            c.createError = 1;
            return;
        }
        $http({
            method: 'POST',
            url: '/create/' + c.name + '/' + c.password,
            data: c.livesplit
        }).then(function successCallback(response) {
            window.location.href = window.location.href + "route.html?g=" + c.name;
        }, function errorCallback(response) {
            c.createError = response.status;
        });
    };
});

app.controller('ListCtrl', function($http) {

    var l = this;

    l.list = [];

    init = function() {
        $http({
            method: 'GET',
            url: '/graphs'
        }).then(function successCallback(response) {
            l.list = response.data;
        }, function errorCallback(response) {
            log("404 I guess...");
        });
    }

    init();
});